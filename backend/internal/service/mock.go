package service

import (
	"log"
)

// StartPriceSimulationMock 启动全量市场模拟
func (s *MarketService) StartPriceSimulationMock(symbol string) {
	log.Printf("市场全量数据模拟启动: %s", symbol)

	// 模拟交易所内部成交
	go s.simulateExchangeTrades(symbol)

	// 模拟 Uniswap 链上事件
	go s.simulateUniswapEvents(symbol)
}
