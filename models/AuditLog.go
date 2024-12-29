package models

import (
	"errors"
	"time"
)

type AuditLog struct {
	ID         uint   `gorm:"primaryKey"`
	EntityType string `gorm:"not null"`
	EntityID   uint   `gorm:"not null;index"`
	Action     string `gorm:"not null"`
	Details    string `gorm:"type:text"`
	CreatedAt  time.Time
}

var validActions = map[string]bool{"create": true, "update": true, "delete": true}
var validEntityTypes = map[string]bool{"user": true, "transaction": true}

func NewAuditLog(entityType, action string, entityID uint, details string) (*AuditLog, error) {
	if !validActions[action] {
		return nil, errors.New("invalid action")
	}
	if !validEntityTypes[entityType] {
		return nil, errors.New("invalid entity type")
	}

	return &AuditLog{
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		Details:    details,
		CreatedAt:  time.Now(),
	}, nil
}
