package models

import "time"

type Rule struct {
	ID                  uint    `gorm:"primaryKey"`
	UserID              uint    `gorm:"not null"`
	User                User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	MaxAmountToTransfer float64 `gorm:"default:0"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
