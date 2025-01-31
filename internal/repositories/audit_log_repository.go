package repositories

import (
	"errors"
	"gorm.io/gorm"
	"ptm/internal/db"
	"ptm/internal/models"
	"ptm/pkg/utils/customError"
)

type AuditLogRepository interface {
	Create(log *models.AuditLog) (*models.AuditLog, error)
	GetModelLogs(entityType string, entityID uint) ([]models.AuditLog, error)
}

type auditLogRepository struct {
}

func NewAuditLogRepository() AuditLogRepository {
	return &auditLogRepository{}
}

func (a *auditLogRepository) Create(log *models.AuditLog) (*models.AuditLog, error) {
	created := db.DB.Create(log)

	if created.Error != nil {
		if errors.Is(created.Error, gorm.ErrRecordNotFound) {
			return nil, customError.NotFound("audit log not found")
		}
		return nil, customError.InternalServerError("something went wrong")
	}
	return log, nil
}

func (a *auditLogRepository) Get(id uint) (*models.AuditLog, error) {
	var log models.AuditLog

	err := db.DB.First(&log, id).Error
	if err != nil {
		return nil, customError.NotFound("audit_log not found")
	}

	return &log, nil
}

func (a *auditLogRepository) GetModelLogs(entityType string, entityID uint) ([]models.AuditLog, error) {
	var logs []models.AuditLog

	err := db.DB.Where("entity_type = ? AND entity_id = ?", entityType, entityID).Find(&logs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customError.NotFound("no logs found for the specified entity")
		}
		return nil, customError.InternalServerError("something went wrong")
	}
	return logs, nil
}

func (a *auditLogRepository) Delete(id uint) error {
	log, err := a.Get(id)

	if err != nil {
		return customError.InternalServerError("something went wrong")
	}

	err = db.DB.Delete(log).Error

	if err != nil {
		return customError.InternalServerError("something went wrong")
	}

	return nil
}
