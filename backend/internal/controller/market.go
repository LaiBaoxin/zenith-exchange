package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/internal/service"
	"github.com/wwater/zenith-exchange/backend/pkg/response"
)

type MarketHandler struct {
	marketService *service.MarketService
	matchService  *service.MatchService
}

func NewMarketHandler(s *service.MarketService, ms *service.MatchService) *MarketHandler {
	return &MarketHandler{
		marketService: s,
		matchService:  ms,
	}
}

// GetKLines 获取 K 线历史数据
func (h *MarketHandler) GetKLines(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "BTC_USDT")
	interval := c.DefaultQuery("interval", "1m")
	limitStr := c.DefaultQuery("limit", "100")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的 limit 参数")
		return
	}

	klines, err := h.marketService.GetKLines(c.Request.Context(), symbol, interval, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取 K 线数据失败")
		return
	}

	response.Success(c, klines)
}

// GetDepth 获取当前盘口深度数据
func (h *MarketHandler) GetDepth(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "BTC_USDT")
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)

	bids, asks := h.matchService.GetDepth(symbol, limit)

	response.Success(c, gin.H{
		"symbol": symbol,
		"bids":   bids,
		"asks":   asks,
	})
}
