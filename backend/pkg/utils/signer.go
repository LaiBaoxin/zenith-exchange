package utils

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wwater/zenith-exchange/backend/pkg/config"
	"os"
	"strings"
)

// LoadPrivateKey 从文件读取私钥
func LoadPrivateKey() (*ecdsa.PrivateKey, error) {
	content, err := os.ReadFile(config.GlobalConfig.Blockchain.KeyPath)
	if err != nil {
		return nil, err
	}
	hexKey := strings.TrimSpace(string(content))
	return crypto.HexToECDSA(hexKey)
}

// SignWithdrawData 生成符合合约校验的签名
func SignWithdrawData(userAddr string, amountStr string) (string, error) {
	privateKey, err := LoadPrivateKey()
	if err != nil {
		return "", err
	}

	// 对数据进行 Keccak256 哈希
	data := []byte(userAddr + amountStr)
	hash := crypto.Keccak256Hash(data)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return "", err
	}

	return common.Bytes2Hex(signature), nil
}
