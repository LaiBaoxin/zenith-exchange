package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/internal/controller"
	"github.com/wwater/zenith-exchange/backend/internal/middleware"
)

// SetupRouter 接收所有需要的 Handler 作为参数
func SetupRouter(vaultH *controller.VaultHandler) *gin.Engine {
	r := gin.Default()

	// 注册全局中间件
	r.Use(CORSMiddleware())

	// 基础健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	{
		vault := api.Group("/vault").Use(middleware.AuthMiddleware())
		{
			// 提现签名接口
			vault.POST("/withdraw-sign", vaultH.HandleWithdraw)
		}
	}

	return r
}

// CORSMiddleware 保持跨域处理
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
