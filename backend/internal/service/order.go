package service

import (
	"context"
	"fmt"
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
func (s *OrderService) CancelOrder(ctx context.Context, userID int64, orderID uint64) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var order model.Order
		// 锁住记录并验证状态 (0: 挂单中, 1: 部分成交)
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ? AND status IN (0, 1)", orderID, userID).
			First(&order).Error; err != nil {
			return fmt.Errorf("订单不可撤销、已成交或不存在")
		}

		// 更新状态为 3 (已撤单)
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

			var refundAsset string
			var refundAmount decimal.Decimal

			if order.Side == "buy" {
				refundAsset = parts[1] // USDT
				refundAmount = unfilledAmount.Mul(decimal.NewFromFloat(order.Price))
			} else {
				refundAsset = parts[0] // BTC
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
func (s *OrderService) GetTodayOrders(ctx context.Context, userID int64, symbol string) ([]model.Order, error) {
	var orders []model.Order
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	query := db.DB.WithContext(ctx).Where("user_id = ? AND created_at >= ?", userID, todayStart)
	if symbol != "" {
		query = query.Where("symbol = ?", symbol)
	}

	err := query.Order("created_at DESC").Find(&orders).Error
	return orders, err
}

// GetUserOrderHistory 获取所有历史订单
func (s *OrderService) GetUserOrderHistory(ctx context.Context, userID int64, symbol string, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := db.DB.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", userID)

	if symbol != "" {
		query = query.Where("symbol = ?", symbol)
	}

	// 统计总数
	query.Count(&total)

	// 分页查询并倒序
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error

	return orders, total, err
}

// CreateOrder 下单逻辑
func (s *OrderService) CreateOrder(ctx context.Context, order *model.Order) error {
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var freezeAsset string
		var freezeAmount decimal.Decimal

		parts := strings.Split(order.Symbol, "_")
		if order.Side == "buy" {
			freezeAsset = parts[1]
			freezeAmount = decimal.NewFromFloat(order.Price).Mul(decimal.NewFromFloat(order.Amount))
		} else {
			freezeAsset = parts[0]
			freezeAmount = decimal.NewFromFloat(order.Amount)
		}

		if err := s.matchService.UpdateBalance(tx, order.UserID, freezeAsset, freezeAmount.Neg(), false, 0, "freeze"); err != nil {
			return fmt.Errorf("余额不足")
		}
		if err := s.matchService.UpdateBalance(tx, order.UserID, freezeAsset, freezeAmount, true, 0, "freeze"); err != nil {
			return err
		}

		return tx.Create(order).Error
	})

	if err == nil {
		go s.matchService.ProcessOrder(order)
	}
	return err
}

// GetOrderDetail 获取订单详情
func (s *OrderService) GetOrderDetail(ctx context.Context, userID uint64, orderID uint64) (*model.Order, error) {
	var order model.Order
	err := db.DB.WithContext(ctx).
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error

	if err != nil {
		return nil, err
	}
	return &order, nil
}

// unfreezeBalance 解冻资产
func (s *OrderService) unfreezeBalance(tx *gorm.DB, userID int64, asset string, amount decimal.Decimal) error {
	if err := s.matchService.UpdateBalance(tx, userID, asset, amount.Neg(), true, 0, "unfreeze"); err != nil {
		return err
	}
	return s.matchService.UpdateBalance(tx, userID, asset, amount, false, 0, "unfreeze")
}
