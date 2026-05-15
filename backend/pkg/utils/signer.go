package utils

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wwater/zenith-exchange/backend/pkg/config"
	"math/big"
	"os"
	"strings"
	"time"
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
func SignWithdrawData(
	userAddr string, // 用户钱包地址 (msg.sender)
	tokenAddr string, // 代币合约地址 (token)
	amountStr string, // 提现金额 (amount)
	nonce int64, // 随机数 (nonce)
	vaultAddr string, // 金库合约地址 (address(this))
	chainID int64, // 链 ID (block.chainid)
) (string, error) {
	// 1. 加载后端私钥 (用于签发授权)
	privateKey, err := LoadPrivateKey()
	if err != nil {
		return "", fmt.Errorf("加载私钥失败: %v", err)
	}

	// 数据类型转换
	userAddress := common.HexToAddress(userAddr)
	tokenAddress := common.HexToAddress(tokenAddr)
	vaultAddress := common.HexToAddress(vaultAddr)

	amount, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		return "", fmt.Errorf("无效的金额格式")
	}

	var data []byte
	data = append(data, userAddress.Bytes()...)                     // msg.sender
	data = append(data, tokenAddress.Bytes()...)                    // token
	data = append(data, common.LeftPadBytes(amount.Bytes(), 32)...) // amount
	data = append(data, math.U256Bytes(big.NewInt(nonce))...)       // nonce
	data = append(data, vaultAddress.Bytes()...)                    // address(this)
	data = append(data, math.U256Bytes(big.NewInt(chainID))...)     // block.chainid

	// 计算 Keccak256 哈希
	hash := crypto.Keccak256Hash(data)

	// 加上以太坊签名消息前缀 (EIP-191)
	// "\x19Ethereum Signed Message:\n32" + hash
	prefixedHash := crypto.Keccak256(
		[]byte("\x19Ethereum Signed Message:\n32"),
		hash.Bytes(),
	)

	// 执行签名
	signature, err := crypto.Sign(prefixedHash, privateKey)
	if err != nil {
		return "", fmt.Errorf("签名失败: %v", err)
	}

	// 调整 V 值
	signature[64] += 27

	return "0x" + common.Bytes2Hex(signature), nil
}

// GenerateNonce 生成一个纳秒级的时间戳作为 Nonce
func GenerateNonce() int64 {
	return time.Now().UnixNano()
}
