package model

import "time"

// BalanceLog 资产变更流水
type BalanceLog struct {
	UserID     uint64    `json:"user_id"`
	Currency   string    `json:"currency"`
	ChangeType string    `json:"change_type"` // deposit, withdraw, trade, freeze
	Amount     float64   `json:"amount"`
	Balance    float64   `json:"balance"`
	LogTime    time.Time `json:"log_time"`
}
