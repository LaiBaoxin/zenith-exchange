package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/internal/controller"
)

func main() {
	r := gin.Default()

	// 路由设置
	r.POST("/api/v1/vault/withdraw-sign", controller.HandleGetWithdrawSignature)

	r.Run(":8888")
}
