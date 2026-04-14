package model

import "time"

// Order 订单表
type Order struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int64     `gorm:"not null;index:idx_user_status" json:"user_id"`
	Symbol       string    `gorm:"type:varchar(20);not null;index:idx_symbol_status" json:"symbol"`
	Side         string    `gorm:"type:enum('buy','sell');not null" json:"side"`
	Type         string    `gorm:"type:enum('limit','market');not null;default:'limit'" json:"type"`
	Price        float64   `gorm:"type:decimal(36,18);not null" json:"price"`
	Amount       float64   `gorm:"type:decimal(36,18);not null" json:"amount"`
	FilledAmount float64   `gorm:"type:decimal(36,18);not null;default:0" json:"filled_amount"`
	Status       int8      `gorm:"not null;default:0;index:idx_user_status;index:idx_symbol_status" json:"status"` // 0:挂单, 1:部分, 2:全额, 3:撤单
	MsgHash      string    `gorm:"type:char(66)" json:"msg_hash"`
	Signature    string    `gorm:"type:text" json:"signature"`
	CreatedAt    time.Time `gorm:"type:datetime(3);default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:datetime(3);default:CURRENT_TIMESTAMP(3)" json:"updated_at"`
	IsMock       bool      `gorm:"type:TINYINT;default:0" json:"is_mock"`
}
