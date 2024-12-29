package models

import (
	"time"
)

type AuditLog struct {
	ID         uint   `gorm:"primaryKey"`
	EntityType string `gorm:"not null"` // Example: "user", "transaction"
	EntityID   uint   `gorm:"not null;index"`
	Action     string `gorm:"not null"` // Example: "create", "update", "delete"
	Details    string `gorm:"type:text"`
	CreatedAt  time.Time
}
