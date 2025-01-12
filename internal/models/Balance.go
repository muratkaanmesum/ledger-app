package models

import (
	"sync"
	"time"
)

type Balance struct {
	UserID        uint    `gorm:"primaryKey"`
	Amount        float64 `gorm:"default:0.0"`
	LastUpdatedAt time.Time
	mutex         sync.RWMutex
	User          User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
