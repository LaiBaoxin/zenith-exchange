package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wwater/zenith-exchange/backend/internal/db"
)

// KlineItem 对应前端常用的 K 线数据结构
type KlineItem struct {
	TS     int64   `json:"ts"`
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

// RunAggregator 启动定时任务
func (s *KlineService) RunAggregator() {
	// 每分钟执行一次基础聚合（1m）
	ticker := time.NewTicker(time.Minute)
	log.Println("K线聚合服务已启动...")

	// 启动时立即尝试补全最近的数据（可选）
	s.AggregateAllPeriods()

	for range ticker.C {
		s.AggregateAllPeriods()
	}
}

// AggregateAllPeriods 聚合多个时间周期
func (s *KlineService) AggregateAllPeriods() {
	periods := []string{"1m", "5m", "15m", "1h", "1d"}
	for _, p := range periods {
		if err := s.AggregateKline(p); err != nil {
			log.Printf("[%s] 聚合失败: %v", p, err)
		}
	}
}

// AggregateKline 通用的聚合逻辑
func (s *KlineService) AggregateKline(period string) error {
	var table string
	var interval string

	// 根据周期映射 ClickHouse 表和时间窗口函数
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

	now := time.Now()
	startTime := now.Add(-time.Hour) // 每次覆盖近一小时的数据，防止漏算

	query := fmt.Sprintf(`
		INSERT INTO %s
		SELECT 
			symbol,
			argMin(price, ts) as open,
			max(price) as high,
			min(price) as low,
			argMax(price, ts) as close,
			sum(amount) as volume,
			%s as ts_window
		FROM trades
		WHERE ts >= ?
		GROUP BY symbol, ts_window
	`, table, interval)

	return db.CH.Exec(context.Background(), query, startTime)
}

// GetKlines 获取 K 线历史数据
func (s *KlineService) GetKlines(ctx context.Context, symbol string, period string, limit int) ([]KlineItem, error) {
	var table string
	switch period {
	case "1m", "5m", "15m", "1h", "1d":
		table = "klines_" + period
	default:
		table = "klines_1m"
	}

	// 按时间倒序
	query := fmt.Sprintf(`
		SELECT 
			toUnixTimestamp(ts) * 1000 as ts_ms,
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
			return nil, err
		}
		list = append(list, item)
	}

	// 返回给前端的数据通常是按时间正序排列的
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}

	return list, nil
}
