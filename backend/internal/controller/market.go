package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/internal/service"
)

type MarketHandler struct {
	marketService *service.MarketService
}

func NewMarketHandler(s *service.MarketService) *MarketHandler {
	return &MarketHandler{marketService: s}
}

// GetKLines 获取 K 线历史数据
func (h *MarketHandler) GetKLines(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "BTC_USDT")
	interval := c.DefaultQuery("interval", "1m") // 1m, 5m, 1h, 1d
	limitStr := c.DefaultQuery("limit", "100")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 limit 参数"})
		return
	}

	// 调用 Service 层从 ClickHouse 查询
	klines, err := h.marketService.GetKLines(c.Request.Context(), symbol, interval, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取 K 线数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": klines,
	})
}
