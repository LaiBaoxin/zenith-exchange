package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wwater/zenith-exchange/backend/internal/contract"
	"github.com/wwater/zenith-exchange/backend/internal/db"
	"github.com/wwater/zenith-exchange/backend/internal/model"
	"github.com/wwater/zenith-exchange/backend/internal/service"
)

type DepositMonitor struct {
	client      *ethclient.Client
	vaultAddr   common.Address
	contractAbi abi.ABI
	hub         *service.Hub
}

// NewDepositMonitor 存款监听器
func NewDepositMonitor(rpcUrl, vaultAddr string, hub *service.Hub) *DepositMonitor {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatalf("RPC连接失败: %v", err)
	}

	parsedAbi, err := abi.JSON(strings.NewReader(contract.ZenithVaultMetaData.ABI))
	if err != nil {
		log.Fatalf("解析合约 ABI 失败: %v", err)
	}

	return &DepositMonitor{
		client:      client,
		vaultAddr:   common.HexToAddress(vaultAddr),
		contractAbi: parsedAbi,
		hub:         hub,
	}
}

// Start 开启监听函数
func (m *DepositMonitor) Start() {
	log.Printf("⚡️ 启动 Deposit 事件监听器 (Vault: %s)...", m.vaultAddr.Hex())

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
			log.Printf("监听错误: %v", err)
		case vLog := <-logs:
			m.handleDepositEvent(vLog)
		}
	}
}

// handleDepositEvent 处理链上日志
func (m *DepositMonitor) handleDepositEvent(vLog types.Log) {
	var event struct {
		Amount   *big.Int
		Currency string
	}

	// 解析非 indexed 数据
	err := m.contractAbi.UnpackIntoInterface(&event, "Deposit", vLog.Data)
	if err != nil {
		log.Printf("解析事件失败: %v", err)
		return
	}

	// 解析 indexed 数据 (user 地址在 Topics[1])
	userAddr := common.HexToAddress(vLog.Topics[1].Hex()).Hex()
	amountStr := event.Amount.String()

	fmt.Printf("监听到充值: 用户=%s, 金额=%s, 币种=%s\n", userAddr, amountStr, event.Currency)

	// 执行数据库逻辑并推送消息
	m.processBalanceAndNotify(userAddr, amountStr, event.Currency)
}

// processBalanceAndNotify 处理余额并通知
func (m *DepositMonitor) processBalanceAndNotify(addr string, amountStr string, currency string) {
	// 查询本地用户 ID
	var user model.User
	if err := db.DB.Where("wallet_address = ?", addr).First(&user).Error; err != nil {
		log.Printf("忽略充值：地址 %s 未在系统注册", addr)
		return
	}

	// 更新 MySQL 余额
	err := db.DB.Model(&model.Account{}).
		Where("user_id = ? AND currency = ?", user.ID, currency).
		UpdateColumn("available", db.DB.Raw("available + ?", amountStr)).Error

	if err != nil {
		log.Printf("数据库更新失败: %v", err)
		return
	}

	// 通过 WebSocket 推送给前端
	msg, _ := json.Marshal(map[string]interface{}{
		"type": "BALANCE_UPDATE",
		"data": map[string]interface{}{
			"currency": currency,
			"amount":   amountStr,
			"message":  "存入成功",
		},
	})

	// 检查用户是否在线并推送
	if client, ok := m.hub.Clients[user.ID]; ok {
		select {
		case client.Send <- msg:
			log.Printf("已向用户 %d 推送余额更新", user.ID)
		default:
			log.Printf("用户 %d 推送通道阻塞", user.ID)
		}
	}
}
