package Api

import (
	"Go_Eth_client/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetBlockInfoByNumber 根据区块号获取区块信息
func GetBlockInfoByNumber(c *gin.Context) {
	// 1. 获取路径/查询参数，这里使用 query 参数 block_num
	numStr := c.Query("block_num")
	if numStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数 block_num 不能为空",
			"data": nil,
		})
		return
	}

	// 2. 参数类型转换
	num, err := strconv.ParseUint(numStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "block_num 必须为合法数字",
			"data": nil,
		})
		return
	}

	// 3. 调用业务层
	blockInfo, err := service.GetBlockInfoByNumber(num)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询区块信息失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	// 4. 成功返回
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": blockInfo,
	})
}

// GetTxInfoByHashHandler 根据交易哈希获取交易详情接口
func GetTxInfoByHashHandler(c *gin.Context) {
	// 1. 从URL参数获取交易哈希
	txHash := c.Query("tx_hash")
	if txHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "tx_hash 参数不能为空",
			"data": nil,
		})
		return
	}

	// 2. 调用业务方法查询数据
	resp, err := service.GetTxByHash(txHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  fmt.Sprintf("查询交易失败: %v", err),
			"data": nil,
		})
		return
	}

	// 3. 正常返回数据
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": resp,
	})
}
