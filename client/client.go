package client

import (
	"Go_Eth_client/config"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// ===================== 全局单例客户端（私有，外部不可直接修改） =====================
var (
	// httpCli 全局ETH HTTP客户端，用于查询区块、余额、合约调用
	httpCli *ethclient.Client
	// wsCli 全局ETH WS客户端，用于事件订阅
	wsCli    *rpc.Client
	erc20ABI abi.ABI

	once    sync.Once // 保证全局初始化只执行一次
	initErr error     // 初始化错误缓存
)

// ===================== 初始化方法（项目启动调用一次） =====================

// InitGlobalClient 统一初始化全局 HTTP + WS 客户端
func InitGlobalClient(AppConfig config.Config) error {
	httpRPC := AppConfig.Eth.HttpRPC
	abiPath := AppConfig.Abi.Path

	// 1. 初始化 HTTP Client
	hc, err := ethclient.Dial(httpRPC)
	if err != nil {
		initErr = fmt.Errorf("http rpc 连接失败: %w", err)
		return err
	}
	httpCli = hc

	// 3. 加载并解析 ERC20 ABI
	abiBytes, err := os.ReadFile(abiPath)
	if err != nil {
		initErr = fmt.Errorf("读取ABI文件失败: %w", err)
		return err
	}
	if err = json.Unmarshal(abiBytes, &erc20ABI); err != nil {
		initErr = fmt.Errorf("解析ABI失败: %w", err)
		return err
	}

	return initErr
}

// ===================== 对外暴露获取方法（只读） =====================

// GetHTTPClient 获取全局HTTP客户端
func GetHTTPClient() *ethclient.Client {
	return httpCli
}

// GetWSClient 获取全局WS客户端
func GetWSClient() *rpc.Client {
	return wsCli
}

// GetERC20ABI 获取全局ERC20 ABI解析器
func GetERC20ABI() abi.ABI {
	return erc20ABI
}

// ===================== 优雅关闭连接（程序退出时调用） =====================

// CloseClient 关闭全局客户端连接
func CloseClient() {
	if httpCli != nil {
		httpCli.Close()
	}
	if wsCli != nil {
		wsCli.Close()
	}
}
