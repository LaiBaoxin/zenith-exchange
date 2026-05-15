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
	assetsH *controller.AssetsHandler,
	wsH *controller.WSHandler,
	marketH *controller.MarketHandler,
	orderH *controller.OrderHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(CORSMiddleware())

	// 公共接口 (无需 Token)
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
		api.POST("/auth/login", authH.Login)
		api.GET("/market/kline", marketH.GetKLines)
		api.GET("/market/depth", marketH.GetDepth)
		api.GET("/system/config", sysH.GetConfig)

		// 特殊：WS 接口单独挂载中间件
		api.GET("/ws", middleware.AuthMiddleware(), wsH.HandleWS)
	}

	// 私有接口
	auth := api.Group("")
	auth.Use(middleware.AuthMiddleware())
	{
		// 资产相关
		assets := auth.Group("/assets")
		{
			assets.GET("/balance", assetsH.GetBalance)
		}

		// 提现相关
		vault := auth.Group("/vault")
		{
			vault.POST("/withdraw-sign", vaultH.HandleWithdraw)
		}

		// 订单相关
		orders := auth.Group("/order")
		{
			orders.GET("/today", orderH.GetTodayList)
			orders.POST("/cancel", orderH.Cancel)
			orders.POST("/place", orderH.Place)
			orders.GET("/list", orderH.GetAllOrders)
			orders.GET("/detail/:id", orderH.GetDetail)
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
