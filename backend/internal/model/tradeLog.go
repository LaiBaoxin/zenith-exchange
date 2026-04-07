package model

import "time"

// TradeLog 成交日志流水
type TradeLog struct {
	TradeID   uint64    `json:"trade_id"`
	OrderID   uint64    `json:"order_id"`
	UserID    uint64    `json:"user_id"`
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
	Amount    float64   `json:"amount"`
	Fee       float64   `json:"fee"`
	Side      string    `json:"side"`
	TradeTime time.Time `json:"trade_time"`
}
