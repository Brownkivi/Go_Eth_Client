package config

import (
	"log"

	"github.com/spf13/viper"
)

// 全局配置实例
var AppConfig Config

// 总配置结构体
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Eth    EthConfig    `mapstructure:"eth"`
	Log    LogConfig    `mapstructure:"log"`
	Abi    ABIConfig    `mapstructure:"abi"`
}

// Web服务配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// 以太坊节点配置
type EthConfig struct {
	HttpRPC       string `mapstructure:"http_rpc"`
	WsRPC         string `mapstructure:"ws_rpc"`
	Erc20Contract string `mapstructure:"erc20_contract"`
}

// 日志配置
type LogConfig struct {
	Level string `mapstructure:"level"`
}

// ABI 配置
type ABIConfig struct {
	Path string `mapstructure:"path"`
}

// Init 初始化配置
func Init() {
	// 配置文件类型、路径、文件名
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config") // 配置文件目录

	// 可选：支持读取环境变量（优先级高于配置文件）
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 反序列化到结构体
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("解析配置失败: %v", err)
	}

	log.Println("配置文件加载完成")
}
