package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/internal/service"
	"github.com/wwater/zenith-exchange/backend/pkg/response"
)

type MarketHandler struct {
	marketService *service.MarketService
}

func NewMarketHandler(svc *service.MarketService) *MarketHandler {
	return &MarketHandler{marketService: svc}
}

// GetKLines 获取K线历史数据
func (h *MarketHandler) GetKLines(c *gin.Context) {
	symbol := c.Query("symbol")
	period := c.DefaultQuery("period", "1m")
	limitStr := c.DefaultQuery("limit", "100")

	if symbol == "" {
		response.Error(c, http.StatusBadRequest, "缺少 symbol 参数")
		return
	}

	limit, _ := strconv.Atoi(limitStr)
	if limit > 1000 {
		limit = 1000
	}

	data, err := h.marketService.GetKLines(c.Request.Context(), symbol, period, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取K线失败: "+err.Error())
		return
	}

	response.Success(c, data)
}

// GetDepth 获取盘口深度
func (h *MarketHandler) GetDepth(c *gin.Context) {
	symbol := c.Query("symbol")
	limitStr := c.DefaultQuery("limit", "20")

	if symbol == "" {
		response.Error(c, http.StatusBadRequest, "缺少 symbol 参数")
		return
	}

	limit, _ := strconv.Atoi(limitStr)

	bids, asks := h.marketService.GetMarketDepth(symbol, limit)

	response.Success(c, gin.H{
		"symbol": symbol,
		"bids":   bids,
		"asks":   asks,
		"ts":     strconv.FormatInt(time.Now().UnixMilli(), 10),
	})
}
