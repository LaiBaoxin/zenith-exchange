package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/internal/controller"
	"github.com/wwater/zenith-exchange/backend/internal/middleware"
)

// SetupRouter 接收所有 Handler 实例
func SetupRouter(
	vaultH *controller.VaultHandler,
	authH *controller.AuthHandler,
	sysH *controller.SystemHandler,
) *gin.Engine {
	r := gin.Default()

	// 全局中间件
	r.Use(CORSMiddleware())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		// 不走鉴权
		api.POST("/auth/login", authH.Login) // 登录

		// 需要鉴权的组
		authGroup := api.Group("/", middleware.AuthMiddleware())
		{
			// 获取全局配置
			authGroup.GET("/system/config", sysH.GetConfig)

			// 提现相关组
			vault := authGroup.Group("/vault")
			{
				vault.POST("/withdraw-sign", vaultH.HandleWithdraw)
			}
		}
	}

	return r
}

// CORSMiddleware 处理跨域请求
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
