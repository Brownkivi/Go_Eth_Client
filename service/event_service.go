package service

import (
	"Go_Eth_client/client"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 配置常量
const (
	WS_RPC_URL    = "ws://127.0.0.1:8545"
	ContractAddr  = "0x5FbDB2315678afecb367f032d93F642f64180aa3"
	RetryInterval = 3 * time.Second
	MaxRetry      = 999
)

// handleContractEvent 处理合约事件
func handleContractEvent(eventLog types.Log) {
	log.Printf("收到原始日志，Topics 数量: %d, Data 长度: %d", len(eventLog.Topics), len(eventLog.Data))

	if len(eventLog.Topics) == 0 && len(eventLog.Data) == 0 {
		log.Println("收到空日志，无事件数据，跳过处理")
		return
	}
	var value *uint64
	abi := client.GetERC20ABI()
	err := abi.UnpackIntoInterface(&value, "Transfer", eventLog.Data)
	if err != nil {
		log.Printf("解码Data失败: %v", err)
		return
	}
	fmt.Println("==================== 收到合约事件 ====================")
	fmt.Printf("区块高度: %d\n", eventLog.BlockNumber)
	fmt.Printf("区块哈希: %s\n", eventLog.BlockHash.Hex())
	fmt.Printf("交易哈希: %s\n", eventLog.TxHash.Hex())
	fmt.Printf("合约地址: %s\n", eventLog.Address.Hex())
	fmt.Printf("Topics: %v\n", eventLog.Topics)
	fmt.Printf("Data(hex): %s\n", (value))
	fmt.Println("=======================================================")
}

// eventLoop 事件消费循环
func eventLoop(logChan chan types.Log, sub ethereum.Subscription, localCtx context.Context) {

	for {
		select {
		case eventLog := <-logChan:
			handleContractEvent(eventLog)

		case err := <-sub.Err():
			log.Printf("订阅通道断开: %v", err)
			return

		case <-localCtx.Done():
			log.Println("当前订阅局部上下文关闭")
			return
		}
	}
}

// StartContractSubscribe 带自动重连的合约事件订阅
func StartContractSubscribe(globalCtx context.Context) {
	go func() {
		retryCount := 0
		contractAddress := common.HexToAddress(ContractAddr)

		for {
			// 全局退出判断
			select {
			case <-globalCtx.Done():
				log.Println("全局退出，停止所有重连")
				return
			default:
			}

			if retryCount >= MaxRetry {
				log.Fatalf("达到最大重连次数 %d，退出", MaxRetry)
				return
			}

			log.Printf("第 %d 次尝试连接 WS 节点", retryCount+1)

			// 使用 ethclient 连接 WS 节点（标准方案）
			client, err := ethclient.DialContext(globalCtx, WS_RPC_URL)
			if err != nil {
				log.Printf("WS 连接失败: %v，%v 后重试", err, RetryInterval)
				retryCount++
				time.Sleep(RetryInterval)
				continue
			}
			log.Println("WS 节点连接成功")

			// 局部上下文
			localCtx, localCancel := context.WithCancel(globalCtx)

			// 标准 FilterQuery：监听指定合约所有事件
			filter := ethereum.FilterQuery{
				Addresses: []common.Address{contractAddress},
			}

			logChan := make(chan types.Log)
			// 标准订阅方法，兼容 Hardhat + 所有版本 ethclient
			sub, err := client.SubscribeFilterLogs(localCtx, filter, logChan)
			if err != nil {
				log.Printf("订阅失败: %v", err)
				localCancel()
				client.Close()
				retryCount++
				time.Sleep(RetryInterval)
				continue
			}

			retryCount = 0
			log.Printf("✅ 合约事件监听已就绪: %s", ContractAddr)

			// 消费事件
			eventLoop(logChan, sub, localCtx)

			// 资源清理
			localCancel()
			sub.Unsubscribe()
			client.Close()

			log.Println("当前连接已断开，准备重连...")
			time.Sleep(RetryInterval)
			retryCount++
		}
	}()
}
