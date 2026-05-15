package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/internal/service"
	"github.com/wwater/zenith-exchange/backend/pkg/response"
	"net/http"
)

type AssetsHandler struct {
	assetsService *service.AssetsService
}

func NewAssetsHandler(svc *service.AssetsService) *AssetsHandler {
	return &AssetsHandler{assetsService: svc}
}

// GetBalance 获取用户所有资产余额
func (h *AssetsHandler) GetBalance(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID, ok := userIDRaw.(int64)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "无效的用户凭证")
		return
	}

	balances, err := h.assetsService.GetUserBalances(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取资产失败")
		return
	}

	response.Success(c, balances)
}
