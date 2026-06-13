package service

import (
	"sync"
	"time"
)

// TransferEvent ERC20 Transfer 事件结构
type TransferEvent struct {
	BlockNumber uint64    `json:"block_number"` // 区块号
	BlockHash   string    `json:"block_hash"`   // 区块哈希
	TxHash      string    `json:"tx_hash"`      // 交易哈希
	From        string    `json:"from"`         // 转出地址
	To          string    `json:"to"`           // 转入地址
	Value       string    `json:"value"`        // 转账金额(wei)
	Time        time.Time `json:"time"`         // 事件本地接收时间
}

// RingBuffer 协程安全环形队列（固定长度）
type RingBuffer struct {
	capacity int
	buf      []*TransferEvent
	head     int
	tail     int
	size     int
	mu       sync.RWMutex
}

// NewRingBuffer 创建环形队列，capacity 最大缓存条数
func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		capacity: capacity,
		buf:      make([]*TransferEvent, capacity),
	}
}

// Push 写入事件，满了自动覆盖最早数据
func (r *RingBuffer) Push(evt *TransferEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.buf[r.tail] = evt
	r.tail = (r.tail + 1) % r.capacity

	if r.size < r.capacity {
		r.size++
	} else {
		// 队列已满，head 后移（淘汰旧数据）
		r.head = (r.head + 1) % r.capacity
	}
}

// GetAll 获取所有缓存事件（按时间从旧到新）
func (r *RingBuffer) GetAll() []*TransferEvent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make([]*TransferEvent, 0, r.size)
	if r.size == 0 {
		return res
	}

	for i := 0; i < r.size; i++ {
		idx := (r.head + i) % r.capacity
		res = append(res, r.buf[idx])
	}
	return res
}

// GetLatest 获取最近 n 条事件
func (r *RingBuffer) GetLatest(n int) []*TransferEvent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if n <= 0 || r.size == 0 {
		return nil
	}
	if n > r.size {
		n = r.size
	}

	res := make([]*TransferEvent, 0, n)
	// 倒序取最新数据
	for i := 0; i < n; i++ {
		idx := (r.tail - 1 - i + r.capacity) % r.capacity
		res = append(res, r.buf[idx])
	}
	return res
}

// Size 当前队列元素数量
func (r *RingBuffer) Size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.size
}
