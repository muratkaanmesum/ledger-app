package models

import (
	"gorm.io/gorm"
	"time"
)

type Event struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	EntityID  string    `gorm:"index" json:"entity_id"`
	Type      string    `json:"type"`
	Payload   string    `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
