package service

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wwater/zenith-exchange/backend/internal/model/resp"
	"math/big"
)

func SignWithdraw(userAddr, tokenAddr string, amount *big.Int, nonce uint64, privateKeyHex string) (*resp.WithdrawResult, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, errors.New("invalid private key")
	}

	// 构造打包数据 (Packed Data)参数:user, token, amount, nonce
	data := append(common.HexToAddress(userAddr).Bytes(), common.HexToAddress(tokenAddr).Bytes()...)
	data = append(data, common.LeftPadBytes(amount.Bytes(), 32)...)
	data = append(data, common.LeftPadBytes(big.NewInt(int64(nonce)).Bytes(), 32)...)

	// Keccak256 哈希
	hash := crypto.Keccak256Hash(data)

	// EIP-191 以太坊签名消息前缀
	prefix := []byte("\x19Ethereum Signed Message:\n32")
	prefixedHash := crypto.Keccak256Hash(append(prefix, hash.Bytes()...))

	// 签名
	sig, err := crypto.Sign(prefixedHash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}
	sig[64] += 27 // 修正 V 值

	return &resp.WithdrawResult{
		Signature: hexutil.Encode(sig),
		Nonce:     nonce,
		Amount:    amount.String(),
		Token:     tokenAddr,
	}, nil
}
