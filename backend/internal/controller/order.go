package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/internal/service"
	"github.com/wwater/zenith-exchange/backend/pkg/response"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: svc}
}

// GetTodayList 获取用户今日交易订单列表
func (h *OrderHandler) GetTodayList(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	symbol := c.Query("symbol")

	orders, err := h.orderService.GetTodayOrders(c.Request.Context(), userID.(uint64), symbol)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取订单列表失败")
		return
	}

	response.Success(c, orders)
}

// Cancel 撤销未成交或部分成交的订单
func (h *OrderHandler) Cancel(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	// 接收参数
	var req struct {
		OrderID string `json:"order_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误：需要有效的 order_id")
		return
	}

	// 转换 ID
	orderID, err := strconv.ParseUint(req.OrderID, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的订单ID格式")
		return
	}

	// 执行撤单逻辑
	err = h.orderService.CancelOrder(c.Request.Context(), userID.(uint64), orderID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, "撤单成功")
}
