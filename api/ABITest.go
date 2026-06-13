package Api

import (
	"Go_Eth_client/erc20"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	// Infura Sepolia RPC / 本地 Hardhat RPC
	RPC = "http://127.0.0.1:8545"
	// ERC20 合约地址
	ContractAddr     = "0x5FbDB2315678afecb367f032d93F642f64180aa3"
	QueryAddr        = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
	senderPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	toAddrStr        = " 0x70997970c51812dc3a010c7d01b50e0d17dc79c8"
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

// newAuth 根据私钥构造 TransactOpts（签名器）
func newAuth(client *ethclient.Client, privateKeyHex string) (*bind.TransactOpts, error) {
	// 解析私钥
	privKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	// 获取钱包地址
	publicKey := privKey.Public().(*ecdsa.PublicKey)
	fromAddr := crypto.PubkeyToAddress(*publicKey)

	// 获取账户 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		return nil, err
	}

	// 获取链 ID（防重放攻击）
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	// 构造签名器
	auth, err := bind.NewKeyedTransactorWithChainID(privKey, chainID)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // 转账 ETH 金额，ERC20 交易填 0
	auth.GasLimit = uint64(300000) // 燃气上限
	auth.GasPrice, err = client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	return auth, nil
}

func transfer() {
	client, err := ethclient.Dial(RPC)
	if err != nil {
		log.Fatalf("连接节点失败: %v", err)
	}
	defer client.Close()

	contractAddr := common.HexToAddress(ContractAddr)
	toAddr := common.HexToAddress(toAddrStr)

	// 3. 初始化 ERC20 合约实例
	erc20Instance, err := erc20.NewErc20(contractAddr, client)
	if err != nil {
		log.Fatalf("初始化合约失败: %v", err)
	}

	// 4. 构造交易签名器
	auth, err := newAuth(client, senderPrivateKey)
	if err != nil {
		log.Fatalf("构造签名器失败: %v", err)
	}
	fromAddr := auth.From
	fmt.Printf("发送账户: %s\n", fromAddr.Hex())

	// ===================== 示例1：查询余额（读调用）=====================
	balance, err := erc20Instance.BalanceOf(&bind.CallOpts{}, fromAddr)
	if err != nil {
		log.Fatalf("查询余额失败: %v", err)
	}
	fmt.Printf("账户余额(原始精度): %s\n", balance.String())

	decimals, _ := erc20Instance.Decimals(&bind.CallOpts{})
	fmt.Printf("代币小数位: %d\n", decimals)

	// ===================== 示例2：执行 transfer 转账（主流程）=====================
	// 转账数量：比如转 1 枚代币，根据 decimals 补 0
	amount := new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))

	// 发送交易
	tx, err := erc20Instance.Transfer(auth, toAddr, amount)
	if err != nil {
		log.Fatalf("发起转账交易失败: %v", err)
	}
	fmt.Printf("转账交易已发送，TxHash: %s\n", tx.Hash().Hex())

	// 等待交易上链确认
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatalf("等待交易确认失败: %v", err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Fatal("交易执行失败，状态码非成功")
	}
	fmt.Println("转账交易上链成功！")
}
