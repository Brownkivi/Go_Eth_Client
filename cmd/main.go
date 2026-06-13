package main

import (
	Api "Go_Eth_client/api"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	// config.Init()
	// var Config = config.AppConfig
	// client.InitGlobalClient(Config)

	// // 全局上下文 用于优雅退出
	// globalCtx, globalCancel := context.WithCancel(context.Background())
	// defer globalCancel()

	// // 1. 启动合约事件订阅（带重连）
	// service.StartContractSubscribe(globalCtx)

	// // 异步启动Gin
	// go startGinServer()

	// // ========== 核心修复：阻塞主进程，监听 Ctrl+C 退出信号 ==========
	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// // 阻塞在这里，直到收到退出信号
	// <-sigChan

	// log.Println("\n收到退出信号，开始优雅关闭服务...")
	// // 关闭全局上下文，通知所有订阅/协程退出
	// globalCancel()

	// // 预留短暂时间让资源释放
	// time.Sleep(500 * time.Millisecond)
	// log.Println("服务已安全退出")

	Api.GetBlance()
}

func startGinServer() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	_ = r.SetTrustedProxies([]string{"127.0.0.1"})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "ok"})
	})

	log.Println("Gin 服务启动 :9090")
	if err := r.Run(":9090"); err != nil {
		log.Fatalf("Gin 启动失败: %v", err)
	}
}
