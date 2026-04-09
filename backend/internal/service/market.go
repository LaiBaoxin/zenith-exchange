package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/wwater/zenith-exchange/backend/internal/db"
)

type MarketService struct {
	klineService *KlineService
	matchService *MatchService
	hub          *Hub
}

func NewMarketService(ks *KlineService, ms *MatchService, h *Hub) *MarketService {
	return &MarketService{
		klineService: ks,
		matchService: ms,
		hub:          h,
	}
}

func (s *MarketService) GetKLines(ctx context.Context, symbol, period string, limit int) ([]KlineItem, error) {
	return s.klineService.GetKlines(ctx, symbol, period, limit)
}

// GetMarketDepth 获取当前盘口深度
func (s *MarketService) GetMarketDepth(symbol string, limit int) (bids, asks []PriceLevel) {
	return s.matchService.GetDepth(symbol, limit)
}

// simulateExchangeTrades 模拟内部撮合成交
func (s *MarketService) simulateExchangeTrades(symbol string) {
	ticker := time.NewTicker(2 * time.Second)
	basePrice := 65000.0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for range ticker.C {
		change := (r.Float64() * 100) - 50
		currentPrice := basePrice + change
		amount := r.Float64() * 1.5
		now := time.Now()
		side := "buy"
		if r.Intn(100) >= 50 {
			side = "sell"
		}

		ctx := context.Background()
		// 写入 trades 表 (K线源)
		_ = db.CH.Exec(ctx, "INSERT INTO trades (symbol, price, amount, taker_side, ts) VALUES (?, ?, ?, ?, ?)",
			symbol, fmt.Sprintf("%.18f", currentPrice), fmt.Sprintf("%.18f", amount), side, now)

		// 写入 trade_logs 表 (适配 DDL)
		_ = db.CH.Exec(ctx, "INSERT INTO trade_logs (symbol, price, amount, side, taker_order_id, maker_order_id, ts) VALUES (?, ?, ?, ?, ?, ?, ?)",
			symbol, fmt.Sprintf("%.18f", currentPrice), fmt.Sprintf("%.18f", amount), side, uint64(r.Intn(100000)), uint64(r.Intn(100000)), now)

		// 广播推送
		tradeUpdate := map[string]interface{}{
			"type": "TRADE_UPDATE",
			"data": map[string]interface{}{
				"symbol": symbol,
				"price":  fmt.Sprintf("%.2f", currentPrice),
				"amount": fmt.Sprintf("%.4f", amount),
				"ts":     now.UnixMilli(),
				"side":   side,
			},
		}
		msg, _ := json.Marshal(tradeUpdate)
		if s.hub != nil {
			s.hub.Broadcast <- msg
		}

		if s.matchService != nil {
			s.matchService.BroadcastDepth(symbol)
		}

		basePrice = currentPrice
	}
}

// simulateUniswapEvents 模拟 Uniswap V3 事件 (适配 DDL)
func (s *MarketService) simulateUniswapEvents(symbol string) {
	ticker := time.NewTicker(8 * time.Second)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for range ticker.C {
		now := time.Now()
		amount0 := r.Float64() * 10.0
		amount1 := amount0 * 65000.0 // 模拟交易额

		err := db.CH.Exec(context.Background(),
			"INSERT INTO uniswap_events (block_number, tx_hash, pool_address, tick, sqrt_price_x96, amount0_in, amount1_out, ts) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			uint64(19000000+r.Intn(1000)),
			fmt.Sprintf("0x%x", r.Uint64()),              // 模拟哈希
			"0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640", // 模拟池地址
			int32(200000+r.Intn(1000)),
			"150239281358912384912384", // 模拟 X96 价格
			fmt.Sprintf("%.18f", amount0),
			fmt.Sprintf("%.18f", amount1),
			now,
		)
		if err != nil {
			log.Printf("Uniswap 事件模拟写入失败: %v", err)
		}
	}
}
