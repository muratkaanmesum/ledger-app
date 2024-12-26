package models

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID     uint    `json:"user_id"`
	FromUserId uint    `json:"from_user_id"`
	ToUserId   uint    `json:"to_user_id"`
	Amount     float64 `json:"amount"`
	Type       string  `json:"type"`
}
