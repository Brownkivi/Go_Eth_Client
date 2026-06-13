package Api

import (
	"Go_Eth_client/erc20"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	// Infura Sepolia RPC / 本地 Hardhat RPC
	RPC = "http://127.0.0.1:8545"
	// ERC20 合约地址
	ContractAddr = "0x5FbDB2315678afecb367f032d93F642f64180aa3"
	QueryAddr    = "0x70997970c51812dc3a010c7d01b50e0d17dc79c8"
)

func GetBlance() uint64 {
	// 1. 连接节点
	client, err := ethclient.Dial(RPC)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 2. 合约地址
	contractAddr := common.HexToAddress(ContractAddr)

	// 3. 创建合约实例（核心）
	erc20Inst, err := erc20.NewErc20(contractAddr, client)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 查余额
	userAddr := common.HexToAddress(QueryAddr)

	balance, err := erc20Inst.BalanceOf(&bind.CallOpts{}, userAddr)
	if err != nil {
		log.Fatal("查询余额失败：", err)
	}

	log.Println("余额：", balance.String())
	return balance.Uint64()
}
