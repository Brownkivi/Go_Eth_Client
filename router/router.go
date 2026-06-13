package router

import (
	Api "Go_Eth_client/api"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	// 强制开启 debug 模式
	gin.SetMode(gin.DebugMode)
	// 初始化gin引擎
	r := gin.Default()

	// 接口总分组
	api := r.Group("/api/eth")
	{
		// 区块相关接口
		blockGroup := api.Group("/block")
		{
			// 根据区块号查询区块信息
			blockGroup.GET("/info", Api.GetBlockInfoByNumber)
			// 交易查询接口
			blockGroup.GET("/tx/info", Api.GetTxInfoByHashHandler)
		}
	}

	return r
}
