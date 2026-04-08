package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wwater/zenith-exchange/backend/internal/db"
	"github.com/wwater/zenith-exchange/backend/internal/model"
	"gorm.io/gorm"
)

// OrderItem 内存中的精简订单，保持高频匹配性能
type OrderItem struct {
	ID        uint64
	UserID    uint64
	Price     decimal.Decimal
	Amount    decimal.Decimal
	Timestamp time.Time
}

type OrderBook struct {
	Symbol string
	Bids   []*OrderItem // 买单 (Price DESC)
	Asks   []*OrderItem // 卖单 (Price ASC)
	mu     sync.Mutex
}

type MatchService struct {
	Books map[string]*OrderBook
	hub   *Hub
	mu    sync.RWMutex
}

func NewMatchService(h *Hub) *MatchService {
	return &MatchService{
		Books: make(map[string]*OrderBook),
		hub:   h,
	}
}

// ProcessOrder 接收并处理新订单
func (s *MatchService) ProcessOrder(order *model.Order) {
	s.mu.Lock()
	book, ok := s.Books[order.Symbol]
	if !ok {
		book = &OrderBook{Symbol: order.Symbol}
		s.Books[order.Symbol] = book
	}
	s.mu.Unlock()

	book.mu.Lock()
	defer book.mu.Unlock()

	// 将 float64 转为 decimal 保证撮合精度
	price := decimal.NewFromFloat(order.Price)
	amount := decimal.NewFromFloat(order.Amount)

	newOrder := &OrderItem{
		ID:        order.ID,
		UserID:    order.UserID,
		Price:     price,
		Amount:    amount,
		Timestamp: order.CreatedAt,
	}

	if order.Side == "buy" {
		s.match(book, newOrder, &book.Asks, true)
	} else {
		s.match(book, newOrder, &book.Bids, false)
	}
}

func (s *MatchService) match(book *OrderBook, taker *OrderItem, makers *[]*OrderItem, isTakerBuy bool) {
	remaining := taker.Amount

	for len(*makers) > 0 && remaining.GreaterThan(decimal.Zero) {
		maker := (*makers)[0]

		// 价格匹配逻辑
		canMatch := false
		if isTakerBuy {
			canMatch = taker.Price.GreaterThanOrEqual(maker.Price)
		} else {
			canMatch = taker.Price.LessThanOrEqual(maker.Price)
		}

		if !canMatch {
			break
		}

		matchedAmount := decimal.Min(remaining, maker.Amount)
		matchPrice := maker.Price

		// 执行数据库更新和资产结算
		err := s.handleTrade(book.Symbol, taker, maker, matchPrice, matchedAmount, isTakerBuy)
		if err != nil {
			log.Printf("撮合事务失败: %v", err)
			break
		}

		remaining = remaining.Sub(matchedAmount)
		maker.Amount = maker.Amount.Sub(matchedAmount)

		if maker.Amount.IsZero() {
			*makers = (*makers)[1:]
		}
	}

	// 剩余部分进入订单簿挂单
	if remaining.GreaterThan(decimal.Zero) {
		taker.Amount = remaining
		s.addToOrderBook(book, taker, isTakerBuy)
	}
}

// handleTrade WS 广播
func (s *MatchService) handleTrade(symbol string, taker, maker *OrderItem, price, amount decimal.Decimal, isTakerBuy bool) error {
	now := time.Now()
	side := "sell"
	if isTakerBuy {
		side = "buy"
	}

	// 计算成交额 (Price * Amount)
	totalQuoteAmount := price.Mul(amount)

	// 事务更新：订单状态和账户余额
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 更新订单状态 (Maker & Taker)
		if err := s.updateOrderStatus(tx, maker.ID, amount); err != nil {
			return err
		}
		if err := s.updateOrderStatus(tx, taker.ID, amount); err != nil {
			return err
		}

		// 更新账户余额
		baseAsset := "BTC"
		quoteAsset := "USDT"

		if isTakerBuy {
			// Taker 买, Maker 卖
			// Taker (买家): 扣除 Quote 冻结 (USDT), 增加 Base 可用 (BTC)
			if err := s.updateBalance(tx, taker.UserID, quoteAsset, totalQuoteAmount.Neg(), true); err != nil {
				return err
			}
			if err := s.updateBalance(tx, taker.UserID, baseAsset, amount, false); err != nil {
				return err
			}

			// Maker (卖家): 扣除 Base 冻结 (BTC), 增加 Quote 可用 (USDT)
			if err := s.updateBalance(tx, maker.UserID, baseAsset, amount.Neg(), true); err != nil {
				return err
			}
			if err := s.updateBalance(tx, maker.UserID, quoteAsset, totalQuoteAmount, false); err != nil {
				return err
			}
		} else {
			// Taker 卖, Maker 买
			// Taker (卖家): 扣除 Base 冻结 (BTC), 增加 Quote 可用 (USDT)
			if err := s.updateBalance(tx, taker.UserID, baseAsset, amount.Neg(), true); err != nil {
				return err
			}
			if err := s.updateBalance(tx, taker.UserID, quoteAsset, totalQuoteAmount, false); err != nil {
				return err
			}

			// Maker (买家): 扣除 Quote 冻结 (USDT), 增加 Base 可用 (BTC)
			if err := s.updateBalance(tx, maker.UserID, quoteAsset, totalQuoteAmount.Neg(), true); err != nil {
				return err
			}
			if err := s.updateBalance(tx, maker.UserID, baseAsset, amount, false); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 异步同步到 ClickHouse 和 WebSocket
	s.syncToSecondarySystems(symbol, price, amount, side, now)
	return nil
}

func (s *MatchService) addToOrderBook(book *OrderBook, order *OrderItem, isBuy bool) {
	if isBuy {
		book.Bids = append(book.Bids, order)
		sort.Slice(book.Bids, func(i, j int) bool {
			if book.Bids[i].Price.Equal(book.Bids[j].Price) {
				return book.Bids[i].Timestamp.Before(book.Bids[j].Timestamp)
			}
			return book.Bids[i].Price.GreaterThan(book.Bids[j].Price)
		})
	} else {
		book.Asks = append(book.Asks, order)
		sort.Slice(book.Asks, func(i, j int) bool {
			if book.Asks[i].Price.Equal(book.Asks[j].Price) {
				return book.Asks[i].Timestamp.Before(book.Asks[j].Timestamp)
			}
			return book.Asks[i].Price.LessThan(book.Asks[j].Price)
		})
	}
}

// updateOrderStatus 更新订单成交额和状态
func (s *MatchService) updateOrderStatus(tx *gorm.DB, orderID uint64, amount decimal.Decimal) error {
	return tx.Model(&model.Order{}).Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"filled_amount": gorm.Expr("filled_amount + ?", amount.InexactFloat64()),
			"status":        gorm.Expr("IF(filled_amount + ? >= amount, 2, 1)", amount.InexactFloat64()),
		}).Error
}

// updateBalance 原子更新用户余额
func (s *MatchService) updateBalance(tx *gorm.DB, userID uint64, asset string, change decimal.Decimal, isFrozen bool) error {
	column := "available"
	if isFrozen {
		column = "frozen"
	}

	result := tx.Model(&model.Account{}).
		Where("user_id = ? AND asset = ?", userID, asset).
		Update(column, gorm.Expr(column+" + ?", change.InexactFloat64()))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("账户不存在: UserID=%d, Asset=%s", userID, asset)
	}
	return nil
}

// syncToSecondarySystems 异步推送
func (s *MatchService) syncToSecondarySystems(symbol string, price, amount decimal.Decimal, side string, now time.Time) {
	go func() {
		// 写入 ClickHouse
		_ = db.CH.Exec(context.Background(),
			"INSERT INTO trades (symbol, price, amount, taker_side, ts) VALUES (?, ?, ?, ?, ?)",
			symbol, price.String(), amount.String(), side, now,
		)

		// 广播成交消息
		tradeMsg, _ := json.Marshal(map[string]interface{}{
			"type": "TRADE_UPDATE",
			"data": map[string]interface{}{
				"symbol": symbol,
				"price":  price.String(),
				"amount": amount.String(),
				"side":   side,
				"ts":     now.UnixMilli(),
			},
		})
		s.hub.Broadcast <- tradeMsg
	}()
}

// RemoveFromBook 从内存订单簿中移除指定订单
func (s *MatchService) RemoveFromBook(symbol string, orderID uint64, side string) {
	s.mu.RLock()
	book, ok := s.Books[symbol]
	s.mu.RUnlock()

	if !ok {
		return
	}

	book.mu.Lock()
	defer book.mu.Unlock()

	if side == "buy" {
		book.Bids = s.removeFromSlice(book.Bids, orderID)
	} else {
		book.Asks = s.removeFromSlice(book.Asks, orderID)
	}
}

// removeFromSlice 从切片中滤除订单
func (s *MatchService) removeFromSlice(slice []*OrderItem, id uint64) []*OrderItem {
	for i, item := range slice {
		if item.ID == id {
			// 从切片中移除元素
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
