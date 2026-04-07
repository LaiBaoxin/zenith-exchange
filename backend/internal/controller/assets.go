package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/internal/db"
	"github.com/wwater/zenith-exchange/backend/internal/model"
	"github.com/wwater/zenith-exchange/backend/pkg/response"
	"net/http"
)

type AssetsHandler struct{}

func NewAssetsHandler() *AssetsHandler {
	return &AssetsHandler{}
}

// GetBalance 获取用户所有资产余额
func (h *AssetsHandler) GetBalance(c *gin.Context) {
	// 从 AuthMiddleware 中获取缓存的 userID
	val, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := val.(uint64)

	var accounts []model.Account
	// 优先从 MySQL 查询，高并发场景后续可在此处加入 Redis 缓存逻辑
	if err := db.DB.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "系统繁忙")
		return
	}
	response.Success(c, gin.H{
		"code": 200,
		"data": accounts,
	})
}
