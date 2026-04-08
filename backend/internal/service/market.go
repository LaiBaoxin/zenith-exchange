package service

import (
	"context"
	"log"
	"time"
)

type MarketService struct {
	klineService *KlineService
	matchService *MatchService
}

func NewMarketService(ks *KlineService, ms *MatchService) *MarketService {
	return &MarketService{
		klineService: ks,
		matchService: ms,
	}
}

// GetKLines 获取历史K线
func (s *MarketService) GetKLines(ctx context.Context, symbol, period string, limit int) ([]KlineItem, error) {
	return s.klineService.GetKlines(ctx, symbol, period, limit)
}

// GetMarketDepth 获取当前盘口深度
func (s *MarketService) GetMarketDepth(symbol string, limit int) (bids, asks []PriceLevel) {
	return s.matchService.GetDepth(symbol, limit)
}

// StartPriceSimulationMock 模拟数据生成
func (s *MarketService) StartPriceSimulationMock(symbol string) {
	log.Printf("启动模拟行情引擎: %s", symbol)
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for range ticker.C {
			// 定时触发一次深度广播，激活 WebSocket 推送
			if s.matchService != nil {
				s.matchService.BroadcastDepth(symbol)
			}
		}
	}()
}
