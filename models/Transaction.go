package models

import (
	"gorm.io/gorm"
)

// Transaction represents a financial transaction for a user
type Transaction struct {
	gorm.Model
	UserID uint    `json:"user_id"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
}
