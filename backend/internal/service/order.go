package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wwater/zenith-exchange/backend/internal/db"
	"github.com/wwater/zenith-exchange/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderService struct {
	matchService *MatchService
}

func NewOrderService(ms *MatchService) *OrderService {
	return &OrderService{matchService: ms}
}

// CancelOrder 撤单逻辑
func (s *OrderService) CancelOrder(ctx context.Context, userID uint64, orderID uint64) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var order model.Order
		// 锁住订单记录，防止正在撮合时撤单
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ? AND status IN (0, 1)", orderID, userID).
			First(&order).Error; err != nil {
			return fmt.Errorf("订单不可撤销、已成交或不存在")
		}

		// 更新数据库状态为 3 (已撤单)
		if err := tx.Model(&order).Update("status", 3).Error; err != nil {
			return err
		}

		// 资产解冻逻辑（未成交金额 = 总量 - 已成交量）
		unfilledAmount := decimal.NewFromFloat(order.Amount).Sub(decimal.NewFromFloat(order.FilledAmount))

		if unfilledAmount.GreaterThan(decimal.Zero) {
			parts := strings.Split(order.Symbol, "_")
			if len(parts) != 2 {
				return fmt.Errorf("无效的交易对格式")
			}
			baseAsset := parts[0]  // BTC
			quoteAsset := parts[1] // USDT

			var refundAsset string
			var refundAmount decimal.Decimal

			if order.Side == "buy" {
				// 买单：冻结的是钱 (USDT)，金额 = 未成交数量 * 挂单价格
				refundAsset = quoteAsset
				refundAmount = unfilledAmount.Mul(decimal.NewFromFloat(order.Price))
			} else {
				// 卖单：冻结的是货 (BTC)，金额 = 未成交数量
				refundAsset = baseAsset
				refundAmount = unfilledAmount
			}

			// 执行解冻：减少冻结，增加可用
			if err := s.unfreezeBalance(tx, userID, refundAsset, refundAmount); err != nil {
				return err
			}
		}

		// 从内存订单簿中移除
		s.matchService.RemoveFromBook(order.Symbol, order.ID, order.Side)

		return nil
	})
}

// GetTodayOrders 获取今日订单列表
func (s *OrderService) GetTodayOrders(ctx context.Context, userID uint64, symbol string) ([]model.Order, error) {
	var orders []model.Order
	// 获取今日零点
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	query := db.DB.WithContext(ctx).Where("user_id = ? AND created_at >= ?", userID, todayStart)
	if symbol != "" {
		query = query.Where("symbol = ?", symbol)
	}

	err := query.Order("created_at DESC").Find(&orders).Error
	return orders, err
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, order *model.Order) error {
	// 开启数据库事务：先落库预扣款，成功后再送入撮合引擎
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 计算需要冻结的金额
		var freezeAsset string
		var freezeAmount decimal.Decimal

		parts := strings.Split(order.Symbol, "_")
		if len(parts) != 2 {
			return fmt.Errorf("无效的交易对格式")
		}
		baseAsset := parts[0]
		quoteAsset := parts[1]

		if order.Side == "buy" {
			// 买单：冻结 USDT = 价格 * 数量
			freezeAsset = quoteAsset
			freezeAmount = decimal.NewFromFloat(order.Price).Mul(decimal.NewFromFloat(order.Amount))
		} else {
			// 卖单：冻结 BTC = 数量
			freezeAsset = baseAsset
			freezeAmount = decimal.NewFromFloat(order.Amount)
		}

		// 执行冻结逻辑, 扣除 Available
		if err := s.matchService.updateBalance(tx, order.UserID, freezeAsset, freezeAmount.Neg(), false, 0, "freeze"); err != nil {
			return fmt.Errorf("余额不足: %v", err)
		}
		// 增加 Frozen
		if err := s.matchService.updateBalance(tx, order.UserID, freezeAsset, freezeAmount, true, 0, "freeze"); err != nil {
			return err
		}

		// 插入订单记录
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatalf("createOrder is err:%v", err.Error())
		return err
	}

	// 只有事务返回 nil之后，才进入内存创建订单
	go s.matchService.ProcessOrder(order)

	return nil
}

// unfreezeBalance 解冻资产
func (s *OrderService) unfreezeBalance(tx *gorm.DB, userID uint64, asset string, amount decimal.Decimal) error {
	// 调用 matchService.updateBalance 保证流水一致性
	// 解冻 = 冻结 -amount, 可用 +amount
	if err := s.matchService.updateBalance(tx, userID, asset, amount.Neg(), true, 0, "unfreeze"); err != nil {
		return err
	}
	return s.matchService.updateBalance(tx, userID, asset, amount, false, 0, "unfreeze")
}
