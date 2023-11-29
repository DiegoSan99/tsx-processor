package entity

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Name         string
	Email        string `gorm:"index"`
	Balance      float64
	Transactions []Transaction `gorm:"foreignKey:AccountID"`
}
