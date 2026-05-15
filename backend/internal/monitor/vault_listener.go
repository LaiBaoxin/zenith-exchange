package monitor

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wwater/zenith-exchange/backend/internal/contract"
	"github.com/wwater/zenith-exchange/backend/internal/db"
	"github.com/wwater/zenith-exchange/backend/internal/model"
	"github.com/wwater/zenith-exchange/backend/internal/service"
	"github.com/wwater/zenith-exchange/backend/pkg/config"
	"gorm.io/gorm"
)

var (
	depositEventSig  = crypto.Keccak256Hash([]byte("Deposit(address,address,uint256)"))
	withdrawEventSig = crypto.Keccak256Hash([]byte("Withdraw(address,address,uint256,uint256)"))
)

type VaultMonitor struct {
	client      *ethclient.Client
	vaultAddr   common.Address
	contractAbi abi.ABI
	hub         *service.Hub
}

func NewVaultMonitor(rpcUrl, vaultAddr string, hub *service.Hub) *VaultMonitor {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatalf("RPC连接失败: %v", err)
	}

	parsedAbi, err := abi.JSON(strings.NewReader(contract.ZenithVaultMetaData.ABI))
	if err != nil {
		log.Fatalf("解析合约 ABI 失败: %v", err)
	}

	return &VaultMonitor{
		client:      client,
		vaultAddr:   common.HexToAddress(vaultAddr),
		contractAbi: parsedAbi,
		hub:         hub,
	}
}

func (m *VaultMonitor) Start() {
	log.Printf("VaultMonitor 启动, 监听 %s", m.vaultAddr.Hex())

	query := ethereum.FilterQuery{
		Addresses: []common.Address{m.vaultAddr},
	}

	logs := make(chan types.Log)
	sub, err := m.client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("订阅失败: %v", err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Printf("监听异常: %v", err)
		case vLog := <-logs:
			switch vLog.Topics[0] {
			case depositEventSig:
				m.handleDepositEvent(vLog)
			case withdrawEventSig:
				m.handleWithdrawEvent(vLog)
			}
		}
	}
}

func (m *VaultMonitor) handleDepositEvent(vLog types.Log) {
	var dataEvent struct {
		Amount *big.Int
	}

	err := m.contractAbi.UnpackIntoInterface(&dataEvent, "Deposit", vLog.Data)
	if err != nil {
		log.Printf("Deposit 事件解析失败: %v", err)
		return
	}

	userAddr := common.HexToAddress(vLog.Topics[1].Hex()).Hex()
	tokenAddr := common.HexToAddress(vLog.Topics[2].Hex()).Hex()

	currency := mapTokenToCurrency(tokenAddr)
	decimalAmount := formatWeiToDecimalStr(dataEvent.Amount, 18)
	log.Printf("充值入账: user=%s, amount=%s", userAddr, decimalAmount)

	m.processDepositAndNotify(userAddr, decimalAmount, currency)
}

func (m *VaultMonitor) handleWithdrawEvent(vLog types.Log) {
	var dataEvent struct {
		Amount *big.Int
		Nonce  *big.Int
	}

	err := m.contractAbi.UnpackIntoInterface(&dataEvent, "Withdraw", vLog.Data)
	if err != nil {
		log.Printf("Withdraw 事件解析失败: %v", err)
		return
	}

	userAddr := common.HexToAddress(vLog.Topics[1].Hex()).Hex()
	tokenAddr := common.HexToAddress(vLog.Topics[2].Hex()).Hex()

	currency := mapTokenToCurrency(tokenAddr)
	decimalAmount := formatWeiToDecimalStr(dataEvent.Amount, 18)
	log.Printf("提现确认: user=%s, amount=%s, nonce=%d", userAddr, decimalAmount, dataEvent.Nonce)

	m.processWithdrawAndNotify(userAddr, decimalAmount, currency)
}

func (m *VaultMonitor) processDepositAndNotify(addr string, amountStr string, currency string) {
	var user model.User
	if err := db.DB.Where("LOWER(wallet_address) = LOWER(?)", addr).First(&user).Error; err != nil {
		log.Printf("充值用户未注册: %s", addr)
		return
	}

	db.DB.Model(&model.Account{}).
		Where("user_id = ? AND currency = ?", user.ID, currency).
		UpdateColumn("available", gorm.Expr("available + ?", amountStr))

	msg, _ := json.Marshal(map[string]interface{}{
		"type": "BALANCE_UPDATE",
		"data": map[string]interface{}{"currency": currency, "amount": amountStr},
	})
	if client, ok := m.hub.Clients[user.ID]; ok {
		client.Send <- msg
	}
}

func (m *VaultMonitor) processWithdrawAndNotify(addr string, amountStr string, currency string) {
	var user model.User
	if err := db.DB.Where("LOWER(wallet_address) = LOWER(?)", addr).First(&user).Error; err != nil {
		log.Printf("提现用户未注册: %s", addr)
		return
	}

	result := db.DB.Model(&model.Account{}).
		Where("user_id = ? AND currency = ?", user.ID, currency).
		UpdateColumn("available", gorm.Expr("available - ?", amountStr))
	if result.Error != nil {
		log.Printf("提现扣款失败: %v", result.Error)
		return
	}
	if result.RowsAffected == 0 {
		log.Printf("提现账户不存在: user=%d, currency=%s", user.ID, currency)
		return
	}

	var account model.Account
	db.DB.Where("user_id = ? AND currency = ?", user.ID, currency).First(&account)

	amountFloat, _ := strconv.ParseFloat(amountStr, 64)
	balanceFloat, _ := strconv.ParseFloat(account.Available, 64)
	db.DB.Create(&model.BalanceLog{
		UserID:     int64(user.ID),
		Currency:   currency,
		ChangeType: "withdraw",
		Amount:     -amountFloat,
		Balance:    balanceFloat,
		LogTime:    time.Now(),
	})

	msg, _ := json.Marshal(map[string]interface{}{
		"type": "BALANCE_UPDATE",
		"data": map[string]interface{}{
			"currency": currency,
			"available": account.Available,
			"frozen":    account.Frozen,
		},
	})
	if client, ok := m.hub.Clients[user.ID]; ok {
		client.Send <- msg
	}
}

func mapTokenToCurrency(tokenAddr string) string {
	cfg := config.GlobalConfig.Blockchain
	if strings.EqualFold(tokenAddr, cfg.TokenAddress) {
		return "USDT"
	}
	return "USDT"
}

func formatWeiToDecimalStr(wei *big.Int, decimals int) string {
	if wei == nil {
		return "0"
	}
	s := wei.String()
	if len(s) <= decimals {
		return "0." + strings.Repeat("0", decimals-len(s)) + s
	}
	loc := len(s) - decimals
	return s[:loc] + "." + s[loc:]
}
