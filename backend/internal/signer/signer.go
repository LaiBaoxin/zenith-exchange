package signer

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
)

// GenerateWithdrawSignature 生成提现的 ECDSA 签名
func GenerateWithdrawSignature(
	privateKey *ecdsa.PrivateKey,
	user common.Address,
	token common.Address,
	amount *big.Int,
	nonce *big.Int,
	vaultAddr common.Address,
	chainID *big.Int,
) ([]byte, error) {

	// 模拟 Solidity 的 abi.encodePacked
	var packedData []byte
	packedData = append(packedData, user.Bytes()...)            // address: 20 bytes
	packedData = append(packedData, token.Bytes()...)           // address: 20 bytes
	packedData = append(packedData, math.U256Bytes(amount)...)  // uint256: 32 bytes
	packedData = append(packedData, math.U256Bytes(nonce)...)   // uint256: 32 bytes
	packedData = append(packedData, vaultAddr.Bytes()...)       // address: 20 bytes
	packedData = append(packedData, math.U256Bytes(chainID)...) // uint256: 32 bytes

	// 第一次 Keccak256 哈希 (对应合约里的 messageHash)
	msgHash := crypto.Keccak256(packedData)

	// 模拟 OpenZeppelin 的 MessageHashUtils.toEthSignedMessageHash()
	// 添加以太坊标准签名头: "\x19Ethereum Signed Message:\n" + 消息长度
	prefix := []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(msgHash)))
	ethSignedMessageHash := crypto.Keccak256(append(prefix, msgHash...))

	// 使用 ECDSA 算法进行签名
	signature, err := crypto.Sign(ethSignedMessageHash, privateKey)
	if err != nil {
		return nil, err
	}

	// 调整v值
	signature[64] += 27

	return signature, nil
}
