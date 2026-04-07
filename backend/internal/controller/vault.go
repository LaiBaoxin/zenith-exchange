package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/internal/service"
	"github.com/wwater/zenith-exchange/backend/pkg/config"
	"github.com/wwater/zenith-exchange/backend/pkg/response"
	"math/big"
	"net/http"
)

type VaultHandler struct {
	VaultAddr string
}

func NewVaultHandler(vaultAddr string) *VaultHandler {
	return &VaultHandler{VaultAddr: vaultAddr}
}

func (h *VaultHandler) HandleWithdraw(c *gin.Context) {
	var req struct {
		Amount string `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid params")
		return
	}

	// 从中间件上下文获取地址
	userAddr, _ := c.Get("user_address")
	amountInt, _ := new(big.Int).SetString(req.Amount, 10)

	// 模拟 Nonce 获取逻辑
	var nonce uint64 = 1

	res, err := service.SignWithdraw(
		userAddr.(string),
		config.GlobalConfig.Blockchain.TokenAddress,
		amountInt,
		nonce,
		config.SignerPrivateKey,
	)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, res)
}
