package models

import (
	"ptm/pkg/utils/customError"
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
var validEntityTypes = map[string]bool{"user": true, "transaction": true, "balance": true}

func NewAuditLog(entityType, action string, entityID uint, details string) (*AuditLog, error) {
	if !validActions[action] {
		return nil, customError.BadRequest("Action is not valid")
	}
	if !validEntityTypes[entityType] {
		return nil, customError.BadRequest("Entity is not valid")
	}

	return &AuditLog{
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		Details:    details,
		CreatedAt:  time.Now(),
	}, nil
}
