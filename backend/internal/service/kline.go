package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wwater/zenith-exchange/backend/internal/db"
)

type KlineItem struct {
	TS     uint32  `json:"t"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

type KlineService struct{}

func NewKlineService() *KlineService {
	return &KlineService{}
}

func (s *KlineService) RunAggregator() {
	ticker := time.NewTicker(10 * time.Second)
	log.Println("K线聚合服务已启动...")

	s.AggregateAllPeriods()

	for range ticker.C {
		s.AggregateAllPeriods()
	}
}

func (s *KlineService) AggregateAllPeriods() {
	periods := []string{"1m", "5m", "15m", "1h", "1d"}
	for _, p := range periods {
		if err := s.AggregateKline(p); err != nil {
			log.Printf("[%s] 聚合失败: %v", p, err)
		}
	}
}

func (s *KlineService) AggregateKline(period string) error {
	var table string
	var interval string

	switch period {
	case "1m":
		table, interval = "klines_1m", "toStartOfMinute(ts)"
	case "5m":
		table, interval = "klines_5m", "toStartOfFiveMinutes(ts)"
	case "15m":
		table, interval = "klines_15m", "toStartOfFifteenMinutes(ts)"
	case "1h":
		table, interval = "klines_1h", "toStartOfHour(ts)"
	case "1d":
		table, interval = "klines_1d", "toStartOfDay(ts)"
	default:
		return fmt.Errorf("不支持的周期: %s", period)
	}

	query := fmt.Sprintf(`
	       INSERT INTO %s (symbol, open, high, low, close, volume, ts)
	       SELECT
	          symbol,
	          argMin(price, ts) as open,
	          max(price) as high,
	          min(price) as low,
	          argMax(price, ts) as close,
	          sum(amount) as volume,
	          %s as ts
	       FROM trades
	       WHERE ts >= ?
	       GROUP BY symbol, ts
	    `, table, interval)

	startTime := time.Now().Add(-2 * time.Hour)
	return db.CH.Exec(context.Background(), query, startTime)
}

func (s *KlineService) GetKlines(ctx context.Context, symbol string, period string, limit int) ([]KlineItem, error) {
	var table string
	switch period {
	case "1m", "5m", "15m", "1h", "1d":
		table = "klines_" + period
	default:
		return nil, fmt.Errorf("不支持的周期: %s", period)
	}

	query := fmt.Sprintf(`
	       SELECT
	          toUnixTimestamp(ts) as ts_sec,
	          open, high, low, close, volume
	       FROM %s
	       WHERE symbol = ?
	       ORDER BY ts DESC
	       LIMIT ?
	    `, table)

	rows, err := db.CH.Query(ctx, query, symbol, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []KlineItem
	for rows.Next() {
		var item KlineItem
		if err := rows.Scan(&item.TS, &item.Open, &item.High, &item.Low, &item.Close, &item.Volume); err != nil {
			log.Printf("Scan error: %v", err)
			return nil, err
		}
		list = append(list, item)
	}

	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}

	// 用 trades 表实时数据补充最新 candle，解决 K 线与当前价格不同步的问题
	if len(list) > 0 {
		now := time.Now()
		var startOfPeriod time.Time
		switch period {
		case "1m":
			startOfPeriod = now.Truncate(time.Minute)
		case "5m":
			startOfPeriod = now.Truncate(5 * time.Minute)
		case "15m":
			startOfPeriod = now.Truncate(15 * time.Minute)
		case "1h":
			startOfPeriod = now.Truncate(time.Hour)
		case "1d":
			startOfPeriod = now.Truncate(24 * time.Hour)
		}

		latestTs := time.Unix(int64(list[len(list)-1].TS), 0)
		if !latestTs.Before(startOfPeriod) {
			liveQuery := `
	            SELECT
	               argMin(price, ts) as open,
	               max(price) as high,
	               min(price) as low,
	               argMax(price, ts) as close,
	               sum(amount) as volume
	            FROM trades
	            WHERE symbol = ? AND ts >= ?
	        `
			var live struct {
				Open   float64
				High   float64
				Low    float64
				Close  float64
				Volume float64
			}
			err := db.CH.QueryRow(ctx, liveQuery, symbol, startOfPeriod).Scan(
				&live.Open, &live.High, &live.Low, &live.Close, &live.Volume,
			)
			if err == nil && live.Close > 0 {
				list[len(list)-1] = KlineItem{
					TS:     uint32(startOfPeriod.Unix()),
					Open:   live.Open,
					High:   live.High,
					Low:    live.Low,
					Close:  live.Close,
					Volume: live.Volume,
				}
			}
		}
	}

	return list, nil
}
