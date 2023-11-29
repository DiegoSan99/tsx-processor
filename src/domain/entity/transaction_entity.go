package entity

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	AccountID  uint64
	Account    Account
	Type       string
	Amount     float64
	IntendedAt time.Time `gorm:"autoCreateTime"`
}
