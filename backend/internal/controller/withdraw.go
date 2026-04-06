package controller

import (
	"github.com/wwater/zenith-exchange/backend/internal/model/request"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/internal/signer"
)

func HandleGetWithdrawSignature(c *gin.Context) {
	var req request.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 加载私钥
	privateKey, _ := crypto.HexToECDSA("abc1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab")

	// 转换数据类型
	amount, _ := new(big.Int).SetString(req.Amount, 10)
	nonce := big.NewInt(req.Nonce)
	chainID := big.NewInt(req.ChainID)

	// 调用我们写好的签名服务
	sig, err := signer.GenerateWithdrawSignature(
		privateKey,
		common.HexToAddress(req.User),
		common.HexToAddress(req.Token),
		amount,
		nonce,
		common.HexToAddress(req.VaultAddr),
		chainID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign"})
		return
	}

	// 返回十六进制签名
	c.JSON(http.StatusOK, gin.H{
		"signature": "0x" + common.Bytes2Hex(sig),
	})
}
