package utils

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/wwater/zenith-exchange/backend/internal/contract"
	"log"
	"strings"
)

// GetVaultABI 获取解析后的 ABI 对象，用于日志监听解析
func GetVaultABI() abi.ABI {
	parsed, err := abi.JSON(strings.NewReader(contract.ZenithVaultABI))
	if err != nil {
		log.Fatalf("解析合约 ABI 失败: %v", err)
	}
	return parsed
}
