package service

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wwater/zenith-exchange/backend/internal/db"
	"github.com/wwater/zenith-exchange/backend/internal/model"
	"gorm.io/gorm"
)

type MockService struct {
	MatchSvc *MatchService
	KlineSvc *KlineService
}

func NewMockService(match *MatchService, kline *KlineService) *MockService {
	return &MockService{
		MatchSvc: match,
		KlineSvc: kline,
	}
}

func (s *MockService) StartPriceSimulationMock(symbol string) {
	log.Printf(">>> 模拟器启动: %s", symbol)

	basePrice := s.getLastTradePrice(symbol)
	if basePrice <= 0 {
		basePrice = 65000.0
	}

	go s.StartDepthInjection(symbol, basePrice)
	go s.startActiveMatchingTask(symbol)
	go s.StartUserOrderSweeper(symbol)
}

func (s *MockService) StartUserOrderSweeper(symbol string) {
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		if s.MatchSvc == nil {
			continue
		}

		s.MatchSvc.mu.Lock()
		book, ok := s.MatchSvc.Books[symbol]
		s.MatchSvc.mu.Unlock()
		if !ok {
			continue
		}

		book.mu.Lock()
		var target *OrderItem

		for i, order := range book.Asks {
			if i >= 5 {
				break
			}
			if !order.IsMock {
				target = order
				book.Asks = append(book.Asks[:i], book.Asks[i+1:]...)
				break
			}
		}

		if target == nil {
			for i, order := range book.Bids {
				if i >= 5 {
					break
				}
				if !order.IsMock {
					target = order
					book.Bids = append(book.Bids[:i], book.Bids[i+1:]...)
					break
				}
			}
		}
		book.mu.Unlock()

		if target != nil {
			log.Printf("[Sweeper] 发现真实订单 ID: %d, 正在执行强制成交结算...", target.ID)

			if err := s.finalizeUserOrder(symbol, target); err != nil {
				log.Printf("[Sweeper] 数据库结算异常: %v", err)
				continue
			}

			s.executeInternalOrder(symbol, target.Price.InexactFloat64(), target.Amount.InexactFloat64(), target.UserID != 999)
			s.MatchSvc.BroadcastDepth(symbol)
		}
	}
}

func (s *MockService) finalizeUserOrder(symbol string, item *OrderItem) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&model.Order{}).
			Where("id = ? AND status IN (0, 1)", item.ID).
			Updates(map[string]interface{}{
				"status":        2,
				"filled_amount": gorm.Expr("amount"),
				"updated_at":    time.Now(),
			})

		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil
		}

		parts := strings.Split(symbol, "_")
		base, quote := parts[0], parts[1]
		totalQuote := item.Price.Mul(item.Amount)

		var order model.Order
		if err := tx.First(&order, item.ID).Error; err != nil {
			return err
		}

		if order.Side == "buy" {
			s.MatchSvc.UpdateBalance(tx, order.UserID, quote, totalQuote.Neg(), true, order.ID, "trade_mock")
			s.MatchSvc.UpdateBalance(tx, order.UserID, base, item.Amount, false, order.ID, "trade_mock")
		} else {
			s.MatchSvc.UpdateBalance(tx, order.UserID, base, item.Amount.Neg(), true, order.ID, "trade_mock")
			s.MatchSvc.UpdateBalance(tx, order.UserID, quote, totalQuote, false, order.ID, "trade_mock")
		}
		return nil
	})
}

func (s *MockService) executeInternalOrder(symbol string, price float64, amount float64, isBuy bool) {
	side := "sell"
	if isBuy {
		side = "buy"
	}
	pDec := decimal.NewFromFloat(price)
	aDec := decimal.NewFromFloat(amount)
	now := time.Now()

	_ = db.CH.Exec(context.Background(), `
	        INSERT INTO trades (symbol, price, amount, taker_side, ts) VALUES (?, ?, ?, ?, ?)
	    `, symbol, pDec.String(), aDec.String(), side, now)

	s.updateKlinesAndPush(symbol, price, amount, now)

	if s.MatchSvc != nil && s.MatchSvc.hub != nil {
		msg, _ := json.Marshal(map[string]interface{}{
			"type": "TRADE_UPDATE",
			"data": map[string]interface{}{
				"symbol": symbol,
				"price":  pDec.StringFixed(2),
				"amount": aDec.StringFixed(4),
				"side":   side,
				"ts":     now.UnixMilli(),
			},
		})
		s.MatchSvc.hub.TopicChan <- TopicMessage{Topic: "trade", Symbol: symbol, Message: msg}
	}
}

func (s *MockService) updateKlinesAndPush(symbol string, price float64, amount float64, now time.Time) {
	minuteTs := now.Truncate(time.Minute)
	query := `
	       INSERT INTO klines_1m (symbol, open, high, low, close, volume, ts)
	       SELECT ?, min(open), max(greatest(high, ?)), min(least(low, ?)), ?, sum(volume + ?), ?
	       FROM (
	          SELECT open, high, low, close, volume FROM klines_1m
	          WHERE symbol = ? AND ts = ?
	          UNION ALL
	          SELECT ? as open, ? as high, ? as low, ? as close, ? as volume
	       )`

	_ = db.CH.Exec(context.Background(), query,
		symbol, price, price, price, amount, minuteTs,
		symbol, minuteTs, price, price, price, price, amount,
	)

	if s.MatchSvc != nil && s.MatchSvc.hub != nil {
		klineData, _ := json.Marshal(map[string]interface{}{
			"type": "KLINE_UPDATE",
			"data": map[string]interface{}{
				"symbol": symbol,
				"t":      minuteTs.Unix(),
				"open":   price, "high": price, "low": price, "close": price, "volume": amount,
			},
		})
		s.MatchSvc.hub.TopicChan <- TopicMessage{Topic: "kline", Symbol: symbol, Message: klineData}
	}
}

func (s *MockService) StartDepthInjection(symbol string, basePrice float64) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ticker := time.NewTicker(2 * time.Second)
	curr := basePrice

	for range ticker.C {
		curr += (r.Float64() - 0.5) * 0.0005 * curr
		if s.MatchSvc != nil {
			s.MatchSvc.ClearMockOrders(symbol)
			for i := 1; i <= 10; i++ {
				offset := float64(i) * (2.0 + r.Float64()*3)
				vol := 0.1 + r.Float64()*1.5
				s.MatchSvc.InjectMockOrder(symbol, decimal.NewFromFloat(curr+offset), decimal.NewFromFloat(vol), "sell")
				s.MatchSvc.InjectMockOrder(symbol, decimal.NewFromFloat(curr-offset), decimal.NewFromFloat(vol), "buy")
			}
			s.MatchSvc.BroadcastDepth(symbol)
		}
	}
}

func (s *MockService) startActiveMatchingTask(symbol string) {
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		s.MatchSvc.mu.Lock()
		book, ok := s.MatchSvc.Books[symbol]
		s.MatchSvc.mu.Unlock()
		if !ok || len(book.Asks) == 0 || len(book.Bids) == 0 {
			continue
		}

		book.mu.Lock()
		if len(book.Asks) == 0 || len(book.Bids) == 0 {
			book.mu.Unlock()
			continue
		}
		midPrice := book.Asks[0].Price.Add(book.Bids[0].Price).Div(decimal.NewFromInt(2))
		tradeAmount := 0.01 + rand.Float64()*0.1
		book.mu.Unlock()

		s.executeInternalOrder(symbol, midPrice.InexactFloat64(), tradeAmount, rand.Float64() > 0.5)
	}
}

func (s *MockService) getLastTradePrice(symbol string) float64 {
	var priceStr string
	query := "SELECT price FROM trades WHERE symbol = ? ORDER BY ts DESC LIMIT 1"
	err := db.CH.QueryRow(context.Background(), query, symbol).Scan(&priceStr)
	if err != nil {
		return 0
	}
	p, _ := decimal.NewFromString(priceStr)
	return p.InexactFloat64()
}
