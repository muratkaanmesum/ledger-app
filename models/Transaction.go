package models

import (
	"time"
)

type Transaction struct {
	ID         uint    `gorm:"primaryKey"`
	FromUserID uint    `gorm:"not null;index"`
	ToUserID   uint    `gorm:"not null;index"`
	Amount     float64 `gorm:"not null"`
	Type       string  `gorm:"not null"` // Example: "debit", "credit"
	Status     string  `gorm:"not null"` // Example: "pending", "completed"
	CreatedAt  time.Time

	FromUser User `gorm:"foreignKey:FromUserID;constraint:OnDelete:CASCADE"`
	ToUser   User `gorm:"foreignKey:ToUserID;constraint:OnDelete:CASCADE"`
}
