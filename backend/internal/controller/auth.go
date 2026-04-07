package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/internal/service"
	"github.com/wwater/zenith-exchange/backend/pkg/response"
	"net/http"
)

type AuthHandler struct {
	authService service.AuthService
}

// Login 登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Address string `json:"address" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	token, walletAddr, err := h.authService.LoginByAddress(req.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.Success(c, gin.H{
		"token":         token,
		"walletAddress": walletAddr,
	})
}
