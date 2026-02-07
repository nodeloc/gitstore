package models

import (
	"time"

	"github.com/google/uuid"
)

// ExchangeRate 汇率模型
type ExchangeRate struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	FromCurrency string    `gorm:"type:varchar(3);not null" json:"from_currency"`
	ToCurrency   string    `gorm:"type:varchar(3);not null" json:"to_currency"`
	Rate         float64   `gorm:"type:decimal(18,8);not null" json:"rate"`
	Source       string    `gorm:"type:varchar(50);default:'exchangerate-api'" json:"source"`
	LastUpdated  time.Time `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"last_updated"`
	CreatedAt    time.Time `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName 指定表名
func (ExchangeRate) TableName() string {
	return "exchange_rates"
}
