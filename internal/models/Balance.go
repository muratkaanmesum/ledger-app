package models

import (
	"time"
)

type Balance struct {
	Id            uint    `gorm:"primaryKey"`
	UserID        uint    `gorm:"not null"`
	Amount        float64 `gorm:"default:0.0"`
	LastUpdatedAt time.Time
	User          User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}
