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

type OrderItem struct {
	ID        uint64
	UserID    uint64
	Price     decimal.Decimal
	Amount    decimal.Decimal
	Timestamp time.Time
}

type OrderBook struct {
	Symbol string
	Bids   []*OrderItem
	Asks   []*OrderItem
	mu     sync.Mutex
}

type PriceLevel struct {
	Price  string `json:"price"`
	Amount string `json:"amount"`
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

	newOrder := &OrderItem{
		ID:        order.ID,
		UserID:    order.UserID,
		Price:     decimal.NewFromFloat(order.Price),
		Amount:    decimal.NewFromFloat(order.Amount),
		Timestamp: order.CreatedAt,
	}

	if order.Side == "buy" {
		s.match(book, newOrder, &book.Asks, true)
	} else {
		s.match(book, newOrder, &book.Bids, false)
	}

	s.BroadcastDepth(book.Symbol)
}

func (s *MatchService) match(book *OrderBook, taker *OrderItem, makers *[]*OrderItem, isTakerBuy bool) {
	remaining := taker.Amount

	for len(*makers) > 0 && remaining.GreaterThan(decimal.Zero) {
		maker := (*makers)[0]
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
		// 传入 Taker 和 Maker 的 ID 用于审计日志
		if err := s.handleTrade(book.Symbol, taker, maker, maker.Price, matchedAmount, isTakerBuy); err != nil {
			log.Printf("撮合事务失败: %v", err)
			break
		}

		remaining = remaining.Sub(matchedAmount)
		maker.Amount = maker.Amount.Sub(matchedAmount)

		if maker.Amount.IsZero() {
			*makers = (*makers)[1:]
		}
	}

	if remaining.GreaterThan(decimal.Zero) {
		taker.Amount = remaining
		s.addToOrderBook(book, taker, isTakerBuy)
	}
}

// handleTrade 处理撮合后的数据库事务
func (s *MatchService) handleTrade(symbol string, taker, maker *OrderItem, price, amount decimal.Decimal, isTakerBuy bool) error {
	now := time.Now()
	side := "sell"
	if isTakerBuy {
		side = "buy"
	}
	totalQuoteAmount := price.Mul(amount)

	// 事务更新：订单状态和账户余额
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.updateOrderStatus(tx, maker.ID, amount); err != nil {
			return err
		}
		if err := s.updateOrderStatus(tx, taker.ID, amount); err != nil {
			return err
		}

		// 使用硬编码设置BTC交易
		baseAsset, quoteAsset := "BTC", "USDT"

		if isTakerBuy {
			// Taker 买: 扣冻结 USDT, 加可用 BTC
			if err := s.UpdateBalance(tx, taker.UserID, quoteAsset, totalQuoteAmount.Neg(), true, taker.ID, "trade"); err != nil {
				return err
			}
			if err := s.UpdateBalance(tx, taker.UserID, baseAsset, amount, false, taker.ID, "trade"); err != nil {
				return err
			}
			// Maker 卖: 扣冻结 BTC, 加可用 USDT
			if err := s.UpdateBalance(tx, maker.UserID, baseAsset, amount.Neg(), true, maker.ID, "trade"); err != nil {
				return err
			}
			if err := s.UpdateBalance(tx, maker.UserID, quoteAsset, totalQuoteAmount, false, maker.ID, "trade"); err != nil {
				return err
			}
		} else {
			// Taker 卖: 扣冻结 BTC, 加可用 USDT
			if err := s.UpdateBalance(tx, taker.UserID, baseAsset, amount.Neg(), true, taker.ID, "trade"); err != nil {
				return err
			}
			if err := s.UpdateBalance(tx, taker.UserID, quoteAsset, totalQuoteAmount, false, taker.ID, "trade"); err != nil {
				return err
			}
			// Maker 买: 扣冻结 USDT, 加可用 BTC
			if err := s.UpdateBalance(tx, maker.UserID, quoteAsset, totalQuoteAmount.Neg(), true, maker.ID, "trade"); err != nil {
				return err
			}
			if err := s.UpdateBalance(tx, maker.UserID, baseAsset, amount, false, maker.ID, "trade"); err != nil {
				return err
			}
		}
		return nil
	})

	if err == nil {
		// 这里增加了 taker.ID 和 maker.ID 的传递
		s.syncToSecondarySystems(symbol, price, amount, side, now, taker.ID, maker.ID)
	}
	return err
}

// UpdateBalance 原子更新用户余额
func (s *MatchService) UpdateBalance(tx *gorm.DB, userID uint64, asset string, change decimal.Decimal, isFrozen bool, refID uint64, changeType string) error {
	var account model.Account
	// 获取当前记录
	if err := tx.Where("user_id = ? AND currency = ?", userID, asset).First(&account).Error; err != nil {
		return fmt.Errorf("账户不存在: %d-%s", userID, asset)
	}

	oldVersion := account.Version
	accAvailable, _ := decimal.NewFromString(account.Available)
	accFrozen, _ := decimal.NewFromString(account.Frozen)

	if isFrozen {
		accFrozen = accFrozen.Add(change)
		account.Frozen = accFrozen.String()
	} else {
		accAvailable = accAvailable.Add(change)
		account.Available = accAvailable.String()
	}

	// 版本号自增
	account.Version += 1

	// 执行更新，带上版本号条件
	result := tx.Model(&account).
		Where("id = ? AND version = ?", account.ID, oldVersion).
		Updates(map[string]interface{}{
			"available": account.Available,
			"frozen":    account.Frozen,
			"version":   account.Version,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("资产更新并发冲突 (乐观锁)")
	}

	// 记录流水
	logEntry := model.BalanceLog{
		UserID:     userID,
		Currency:   asset,
		ChangeType: changeType,
		Amount:     change.InexactFloat64(),
		Balance:    accAvailable.InexactFloat64(),
		LogTime:    time.Now(),
	}
	return tx.Create(&logEntry).Error
}

func (s *MatchService) syncToSecondarySystems(symbol string, price, amount decimal.Decimal, side string, now time.Time, takerID, makerID uint64) {
	go func() {
		ctx := context.Background()
		// 写入 trades 表 (K线聚合源)
		_ = db.CH.Exec(ctx, "INSERT INTO trades (symbol, price, amount, taker_side, ts) VALUES (?, ?, ?, ?, ?)",
			symbol, price.String(), amount.String(), side, now)

		// 写入 trade_logs 表
		_ = db.CH.Exec(ctx, "INSERT INTO trade_logs (symbol, price, amount, side, taker_order_id, maker_order_id, ts) VALUES (?, ?, ?, ?, ?, ?, ?)",
			symbol, price.String(), amount.String(), side, takerID, makerID, now)

		// WebSocket
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

		s.hub.TopicChan <- TopicMessage{
			Topic:   "trade",
			Symbol:  symbol,
			Message: tradeMsg,
		}
	}()
}

func (s *MatchService) BroadcastDepth(symbol string) {
	go func() {
		bids, asks := s.GetDepth(symbol, 20)
		msg, _ := json.Marshal(map[string]interface{}{
			"type": "DEPTH_UPDATE",
			"data": map[string]interface{}{
				"symbol": symbol,
				"bids":   bids,
				"asks":   asks,
			},
		})

		s.hub.TopicChan <- TopicMessage{
			Topic:   "depth",
			Symbol:  symbol,
			Message: msg,
		}
	}()
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

func (s *MatchService) GetDepth(symbol string, limit int) (bids, asks []PriceLevel) {
	s.mu.RLock()
	book, ok := s.Books[symbol]
	s.mu.RUnlock()
	if !ok {
		return []PriceLevel{}, []PriceLevel{}
	}

	book.mu.Lock()
	defer book.mu.Unlock()
	bids = aggregateDepth(book.Bids, limit)
	asks = aggregateDepth(book.Asks, limit)
	return bids, asks
}

func aggregateDepth(items []*OrderItem, limit int) []PriceLevel {
	var levels []PriceLevel
	if len(items) == 0 {
		return levels
	}
	var currP, currA decimal.Decimal
	for _, item := range items {
		if currP.IsZero() {
			currP, currA = item.Price, item.Amount
		} else if item.Price.Equal(currP) {
			currA = currA.Add(item.Amount)
		} else {
			levels = append(levels, PriceLevel{Price: currP.String(), Amount: currA.String()})
			if len(levels) >= limit {
				return levels
			}
			currP, currA = item.Price, item.Amount
		}
	}
	if len(levels) < limit && !currP.IsZero() {
		levels = append(levels, PriceLevel{Price: currP.String(), Amount: currA.String()})
	}
	return levels
}

func (s *MatchService) updateOrderStatus(tx *gorm.DB, orderID uint64, amount decimal.Decimal) error {
	return tx.Model(&model.Order{}).Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"filled_amount": gorm.Expr("filled_amount + ?", amount.InexactFloat64()),
			"status":        gorm.Expr("IF(filled_amount + ? >= amount, 2, 1)", amount.InexactFloat64()),
		}).Error
}

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
	s.BroadcastDepth(symbol)
}

func (s *MatchService) removeFromSlice(slice []*OrderItem, id uint64) []*OrderItem {
	for i, item := range slice {
		if item.ID == id {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func (s *MatchService) InitOrderBook() error {
	var activeOrders []model.Order
	if err := db.DB.Where("status IN ?", []int8{0, 1}).Order("created_at ASC").Find(&activeOrders).Error; err != nil {
		return err
	}
	for _, order := range activeOrders {
		s.recoveryOrder(&order)
	}
	return nil
}

func (s *MatchService) recoveryOrder(order *model.Order) {
	s.mu.Lock()
	book, ok := s.Books[order.Symbol]
	if !ok {
		book = &OrderBook{Symbol: order.Symbol}
		s.Books[order.Symbol] = book
	}
	s.mu.Unlock()
	rem := decimal.NewFromFloat(order.Amount).Sub(decimal.NewFromFloat(order.FilledAmount))
	if rem.GreaterThan(decimal.Zero) {
		item := &OrderItem{ID: order.ID, UserID: order.UserID, Price: decimal.NewFromFloat(order.Price), Amount: rem, Timestamp: order.CreatedAt}
		book.mu.Lock()
		s.addToOrderBook(book, item, order.Side == "buy")
		book.mu.Unlock()
	}
}
