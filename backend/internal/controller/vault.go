package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/internal/service"
	"github.com/wwater/zenith-exchange/backend/pkg/config"
	"github.com/wwater/zenith-exchange/backend/pkg/response"
)

type VaultHandler struct {
	vaultAddr    string
	vaultService *service.VaultService
}

func NewVaultHandler(addr string, svc *service.VaultService) *VaultHandler {
	return &VaultHandler{
		vaultAddr:    addr,
		vaultService: svc,
	}
}

func (h *VaultHandler) HandleWithdraw(c *gin.Context) {
	var req struct {
		Amount   string `json:"amount" binding:"required"`
		Currency string `json:"currency" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数格式错误")
		return
	}

	// 安全取值：user_id 和 user_address 由 JWT 中间件注入
	val, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}
	userID, ok := val.(int64)
	if !ok || userID == 0 {
		response.Error(c, http.StatusUnauthorized, "无效的用户凭证")
		return
	}

	walletAddr := c.GetString("user_address")
	if walletAddr == "" {
		response.Error(c, http.StatusUnauthorized, "未能获取钱包地址")
		return
	}

	// 调用 Service 层处理业务逻辑
	sig, nonce, err := h.vaultService.PrepareWithdraw(userID, walletAddr, req.Currency, req.Amount)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 返回前端所需的所有提现凭证 (token 从配置读取，不再硬编码)
	response.Success(c, gin.H{
		"signature": sig,
		"nonce":     nonce,
		"amount":    req.Amount,
		"vault":     h.vaultAddr,
		"token":     config.GlobalConfig.Blockchain.TokenAddress,
	})
}
