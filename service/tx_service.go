package service

import (
	"Go_Eth_client/client"
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TxDetail 交易基础信息
type TxDetail struct {
	From     string `json:"from"`      // 发送地址
	To       string `json:"to"`        // 接收地址
	Value    string `json:"value"`     // 转账金额(wei)
	GasLimit uint64 `json:"gas_limit"` // 交易Gas上限
	Input    string `json:"input"`     // 调用数据/合约入参
	TxHash   string `json:"tx_hash"`   // 交易哈希
}

// TxReceipt 交易回执信息
type TxReceipt struct {
	Status    uint64 `json:"status"`     // 交易状态 1成功 0失败
	GasUsed   uint64 `json:"gas_used"`   // 该交易实际消耗Gas
	LogsCount int    `json:"logs_count"` // 日志条数
}

// TxFullResp 交易完整返回体
type TxFullResp struct {
	Tx      TxDetail  `json:"tx"`
	Receipt TxReceipt `json:"receipt"`
}

// GetTxByHash 根据交易哈希查询 交易 + 回执
func GetTxByHash(txHashStr string) (*TxFullResp, error) {
	txHash := common.HexToHash(txHashStr)
	cli := client.GetHTTPClient()

	// 1. 查询交易本体
	tx, _, err := cli.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("get tx failed: %w", err)
	}

	// 2. 查询交易回执
	receipt, err := cli.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("get receipt failed: %w", err)
	}

	// 3. 解析发送方地址
	chainID, err := cli.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get chain id failed: %w", err)
	}
	// 根据链ID创建对应签名器
	signer := types.LatestSignerForChainID(chainID)
	from, err := types.Sender(signer, tx)
	if err != nil {
		return nil, fmt.Errorf("parse sender failed: %w", err)
	}

	// 处理接收地址（合约创建交易 to 为 nil）
	to := ""
	if tx.To() != nil {
		to = tx.To().Hex()
	}

	// 组装返回数据
	txDetail := TxDetail{
		From:     from.Hex(),
		To:       to,
		Value:    tx.Value().String(),
		GasLimit: tx.Gas(),
		Input:    common.Bytes2Hex(tx.Data()),
		TxHash:   txHashStr,
	}

	receiptInfo := TxReceipt{
		Status:    receipt.Status,
		GasUsed:   receipt.GasUsed,
		LogsCount: len(receipt.Logs),
	}

	return &TxFullResp{
		Tx:      txDetail,
		Receipt: receiptInfo,
	}, nil
}
