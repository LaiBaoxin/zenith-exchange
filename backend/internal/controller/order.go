package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/internal/model"
	"github.com/wwater/zenith-exchange/backend/internal/service"
	"github.com/wwater/zenith-exchange/backend/pkg/response"
	"net/http"
	"strconv"
	"time"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: svc}
}

// GetTodayList 获取用户今日交易订单列表
func (h *OrderHandler) GetTodayList(c *gin.Context) {

	symbol := c.Query("symbol")

	val, _ := c.Get("user_id")
	userID, ok := val.(int64)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "无效的用户 ID 类型")
		return
	}
	orders, err := h.orderService.GetTodayOrders(c.Request.Context(), userID, symbol)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取订单列表失败")
		return
	}

	response.Success(c, orders)
}

// GetAllOrders 获取用户所有历史订单
func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	symbol := c.Query("symbol")

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	orders, total, err := h.orderService.GetUserOrderHistory(c.Request.Context(), userID.(int64), symbol, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取历史账本失败")
		return
	}

	// 返回带 Total 的结构
	response.Success(c, gin.H{
		"list":  orders,
		"total": total,
	})
}

// Cancel 撤销未成交或部分成交的订单
func (h *OrderHandler) Cancel(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	// 支持从 URL 参数获取 ID (符合 DELETE /order/:id 规范)
	idStr := c.Param("id")
	if idStr == "" {
		// 兼容 JSON 提交格式
		var req struct {
			OrderID string `json:"order_id"`
		}
		if err := c.ShouldBindJSON(&req); err == nil {
			idStr = req.OrderID
		}
	}

	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的订单ID")
		return
	}

	err = h.orderService.CancelOrder(c.Request.Context(), userID.(int64), orderID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, "撤单成功")
}

// Place 下单 (保持原有逻辑)
func (h *OrderHandler) Place(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Symbol string  `json:"symbol" binding:"required"`
		Side   string  `json:"side" binding:"required,oneof=buy sell"`
		Price  float64 `json:"price" binding:"required,gt=0"`
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	newOrder := &model.Order{
		UserID:    userID.(int64),
		Symbol:    req.Symbol,
		Side:      req.Side,
		Price:     req.Price,
		Amount:    req.Amount,
		Status:    0,
		CreatedAt: time.Now(),
	}

	if err := h.orderService.CreateOrder(c.Request.Context(), newOrder); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"order_id": strconv.FormatUint(newOrder.ID, 10),
	})
}

// GetDetail 获取单个订单详情
func (h *OrderHandler) GetDetail(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的订单ID")
		return
	}

	order, err := h.orderService.GetOrderDetail(c.Request.Context(), userID.(uint64), orderID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "订单不存在或无权查看")
		return
	}

	response.Success(c, order)
}
