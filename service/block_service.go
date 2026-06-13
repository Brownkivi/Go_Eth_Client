package service

import (
	"Go_Eth_client/client"
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

// BlockInfo 定义返回的区块简要信息结构体
type BlockInfo struct {
	BlockNumber uint64 `json:"block_number"` // 区块高度
	BlockHash   string `json:"block_hash"`   // 区块哈希
	ParentHash  string `json:"parent_hash"`  // 父区块哈希
	Timestamp   uint64 `json:"timestamp"`    // 区块时间戳
	TxCount     int    `json:"tx_count"`     // 交易数量
	GasLimit    uint64 `json:"gas_limit"`    // 区块Gas上限
	GasUsed     uint64 `json:"gas_used"`     // 区块已使用Gas
}

// GetBlockInfoByNumber 根据区块号查询区块信息
func GetBlockInfoByNumber(num uint64) (*BlockInfo, error) {
	client := client.GetHTTPClient()
	// 转为 *big.Int 调用底层API
	blocknum := big.NewInt(int64(num))
	block, err := client.BlockByNumber(context.Background(), blocknum)
	if err != nil {
		return nil, err
	}
	return parseBlock(block), nil
}

// parseBlock 解析区块对象，组装成统一返回结构体
func parseBlock(block *types.Block) *BlockInfo {
	return &BlockInfo{
		BlockNumber: block.NumberU64(),
		BlockHash:   block.Hash().Hex(),
		ParentHash:  block.ParentHash().Hex(),
		Timestamp:   block.Time(),
		TxCount:     len(block.Transactions()),
		GasLimit:    block.GasLimit(),
		GasUsed:     block.GasUsed(),
	}
}
