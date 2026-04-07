package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/wwater/zenith-exchange/backend/internal/db"
	"github.com/wwater/zenith-exchange/backend/internal/model/resp"
)

type MarketService struct {
	hub *Hub
}

func NewMarketService(h *Hub) *MarketService {
	return &MarketService{hub: h}
}

// GetKLines 从 ClickHouse 获取聚合后的 K 线数据
func (s *MarketService) GetKLines(ctx context.Context, symbol string, interval string, limit int) ([]resp.Kline, error) {
	var klines []resp.Kline

	// 使用 ClickHouse 的 toStartOfInterval 进行高效聚合
	query := fmt.Sprintf(`
       SELECT 
          toUnixTimestamp(toStartOfInterval(ts, INTERVAL %s)) as time,
          argMin(price, ts) as open,
          max(price) as high,
          min(price) as low,
          argMax(price, ts) as close,
          cast(sum(amount), 'Float64') as volume
       FROM trades 
       WHERE symbol = ? 
       GROUP BY time 
       ORDER BY time DESC 
       LIMIT ?`, interval)

	rows, err := db.CH.Query(ctx, query, symbol, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var k resp.Kline
		if err := rows.Scan(&k.Time, &k.Open, &k.High, &k.Low, &k.Close, &k.Volume); err != nil {
			return nil, err
		}
		klines = append(klines, k)
	}
	return klines, nil
}

// StartPriceSimulationMock 启动市场价格模拟
func (s *MarketService) StartPriceSimulationMock(symbol string) {
	ticker := time.NewTicker(2 * time.Second)
	basePrice := 65000.0 // 初始价格

	log.Printf("市场模拟引擎已启动 (双向成交模式): %s", symbol)

	go func() {
		// 设置随机数种子，确保每次运行波动不同
		rand.Seed(time.Now().UnixNano())

		for range ticker.C {
			// 产生 [-50.0, 50.0) 的随机波动
			change := (rand.Float64() * 100) - 50
			currentPrice := basePrice + change
			if currentPrice < 0 {
				currentPrice = 1.0
			}

			amount := rand.Float64() * 1.5
			now := time.Now()

			// 随机决定是买还是卖 (TakerSide)
			side := "buy"
			if rand.Intn(100) >= 60 {
				side = "sell"
			}

			// 持久化到 ClickHouse, taker_side 是 Enum8('buy'=1, 'sell'=2)
			err := db.CH.Exec(context.Background(),
				"INSERT INTO trades (symbol, price, amount, taker_side, ts) VALUES (?, ?, ?, ?, ?)",
				symbol,
				fmt.Sprintf("%.18f", currentPrice),
				fmt.Sprintf("%.18f", amount),
				side, // 这里传入随机生成的 side
				now,
			)
			if err != nil {
				log.Printf("ClickHouse 写入模拟数据失败: %v", err)
				continue
			}

			// 构造推送给前端的消息对象
			tradeUpdate := map[string]interface{}{
				"type": "TRADE_UPDATE",
				"data": map[string]interface{}{
					"symbol":    symbol,
					"price":     fmt.Sprintf("%.2f", currentPrice),
					"amount":    fmt.Sprintf("%.4f", amount),
					"ts":        now.UnixMilli(),
					"side":      side, // 告诉前端是买还是卖
					"direction": s.getDirection(change),
				},
			}

			msg, _ := json.Marshal(tradeUpdate)

			// 通过 Hub 进行广播
			s.hub.Broadcast <- msg

			// 更新基准价格
			basePrice = currentPrice
		}
	}()
}

// getDirection 判断价格走势方向
func (s *MarketService) getDirection(change float64) string {
	if change >= 0 {
		return "up"
	}
	return "down"
}
