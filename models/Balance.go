package models

import (
	"gorm.io/gorm"
	"time"
)

type Balance struct {
	gorm.Model
	UserId        uint      `json:"user_id"`
	Amount        uint      `json:"amount"`
	LastUpdatedAt time.Time `gorm:"column:updated_at"`
}
