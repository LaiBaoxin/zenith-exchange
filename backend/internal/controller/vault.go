package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/pkg/response"
	"github.com/wwater/zenith-exchange/backend/pkg/utils"
	"net/http"
)

type VaultHandler struct {
	vaultAddr string
}

func NewVaultHandler(addr string) *VaultHandler {
	return &VaultHandler{vaultAddr: addr}
}

// HandleWithdraw 生成提现签名
func (h *VaultHandler) HandleWithdraw(c *gin.Context) {
	var req struct {
		Amount   string `json:"amount" binding:"required"`
		Currency string `json:"currency" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 获取当前用户钱包地址（从鉴权中间件传过来）
	walletAddr := c.GetString("wallet_address")

	sig, err := utils.SignWithdrawData(walletAddr, req.Amount)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "签名生成失败")
		return
	}
	response.Success(c, gin.H{
		"signature": sig,
		"vault":     h.vaultAddr,
	})
}
