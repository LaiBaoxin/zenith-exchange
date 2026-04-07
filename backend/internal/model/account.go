package model

import "time"

type Account struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"index:uk_user_currency" json:"user_id"`
	Currency  string    `gorm:"size:20;index:uk_user_currency" json:"currency"`
	Available string    `gorm:"type:decimal(36,18);default:'0'" json:"available"`
	Frozen    string    `gorm:"type:decimal(36,18);default:'0'" json:"frozen"`
	Version   uint32    `gorm:"default:0" json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}
