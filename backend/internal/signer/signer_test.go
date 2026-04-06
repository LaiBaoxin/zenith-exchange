package signer

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// TestGenerateWithdrawSignature 测试生产的密钥
func TestGenerateWithdrawSignature(t *testing.T) {
	// 模拟生成一个后端私钥 (实际业务中应从配置或 KMS 加载)
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// 打印出公钥地址，这就相当于合约里的 backendSigner
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	signerAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	t.Logf("Signer Address: %s", signerAddress.Hex())

	// 模拟业务数据
	user := common.HexToAddress("0x1234567890123456789012345678901234567890")
	token := common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")

	// 模拟 50 ether (50 * 10^18)
	amount, _ := new(big.Int).SetString("50000000000000000000", 10)
	nonce := big.NewInt(0)
	vaultAddr := common.HexToAddress("0x9999999999999999999999999999999999999999")
	chainID := big.NewInt(31337) // 模拟 Foundry 本地链 ID

	// 调用我们编写的签名函数
	signature, err := GenerateWithdrawSignature(privateKey, user, token, amount, nonce, vaultAddr, chainID)
	if err != nil {
		t.Fatalf("Failed to generate signature: %v", err)
	}

	// 验证签名长度是否为 65 字节 (r=32, s=32, v=1)
	if len(signature) != 65 {
		t.Errorf("Expected signature length to be 65, got %d", len(signature))
	}

	// 打印出 Hex 格式的签名，通常前缀会加 "0x" 传给前端
	t.Logf("Generated Signature: 0x%s", hex.EncodeToString(signature))

	// 验证 v 值是否已被修正为 27 或 28
	v := signature[64]
	if v != 27 && v != 28 {
		t.Errorf("Expected v to be 27 or 28, got %d", v)
	}
}
