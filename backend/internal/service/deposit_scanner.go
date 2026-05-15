package service

import (
	"context"
	dao "github.com/wwater/zenith-exchange/backend/internal/db"
	"github.com/wwater/zenith-exchange/backend/internal/model"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

type DepositScanner struct {
	Client    *ethclient.Client
	DB        *gorm.DB
	VaultAddr common.Address
	TokenAddr common.Address
}

func NewDepositScanner(rpcURL string, vaultAddr string, tokenAddr string) *DepositScanner {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("无法连接到 Anvil 节点: %v", err)
	}

	return &DepositScanner{
		Client:    client,
		VaultAddr: common.HexToAddress(vaultAddr),
		TokenAddr: common.HexToAddress(tokenAddr),
	}
}

// Start 启动监听协程
func (s *DepositScanner) Start() {
	ctx := context.Background()
	transferSig := []byte("Transfer(address,address,uint256)")
	transferEventHash := crypto.Keccak256Hash(transferSig)

	// 只监听目标代币合约
	query := ethereum.FilterQuery{
		Addresses: []common.Address{s.TokenAddr},
	}

	logs := make(chan types.Log)
	sub, err := s.Client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("正在监听链上充值事件...")

	for {
		select {
		case err := <-sub.Err():
			log.Printf("订阅错误: %v", err)
		case vLog := <-logs:
			if len(vLog.Topics) == 3 && vLog.Topics[0] == transferEventHash {
				from := common.HexToAddress(vLog.Topics[1].Hex())
				to := common.HexToAddress(vLog.Topics[2].Hex())

				// 转给交易所金库 (Vault) 的钱
				if to.Hex() == s.VaultAddr.Hex() {
					amount := new(big.Int).SetBytes(vLog.Data)
					log.Printf("收到充值! From: %s, To Vault, Amount: %s", from.Hex(), amount.String())

					// 更新数据库账本
					s.updateUserBalance(from.Hex(), amount)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

// updateUserBalance 根据转账来源地址给用户加钱
func (s *DepositScanner) updateUserBalance(walletAddress string, amount *big.Int) {
	var users model.User
	err := dao.DB.Where("wallet_address = ?", walletAddress).First(&users).Error
	if err != nil {
		log.Printf("get wallet_address = %s is err: %v", walletAddress, err)
		return
	}

	// 将 big.Int 转为数据库能存的 18 位 Decimal 字符串
	amountStr := formatAmountToDecimal(amount)

	// 执行加钱操作
	err = dao.DB.Exec(`
		UPDATE accounts 
		SET available = available + ? 
		WHERE user_id = ? AND currency = 'USDT'
	`, amountStr, users.ID).Error

	if err != nil {
		log.Printf("充值入账失败: %v", err)
	} else {
		log.Printf("用户 %d 充值入账成功, 增加: %s USDT", users.ID, amountStr)
	}
}

// formatAmountToDecimal 把大数转为带有18位小数的字符串
func formatAmountToDecimal(amount *big.Int) string {
	str := amount.String()
	if len(str) <= 18 {
		pad := strings.Repeat("0", 18-len(str))
		return "0." + pad + str
	}
	return str[:len(str)-18] + "." + str[len(str)-18:]
}
