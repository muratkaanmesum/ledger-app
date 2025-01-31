package services

import (
	"ptm/internal/models"
	"ptm/internal/repositories"
)

type AuditLogService interface {
	CreateLog(entityType, action string, entityId uint) (*models.AuditLog, error)
	FindById(id uint) (*models.AuditLog, error)
	GetModelLogs(entityType string, entityId uint) ([]models.AuditLog, error)
}

type auditLogService struct {
	repo repositories.AuditLogRepository
}

func NewAuditLogService(repository repositories.AuditLogRepository) AuditLogService {
	return &auditLogService{
		repo: repository,
	}
}

func (a *auditLogService) CreateLog(entityType, action string, entityId uint) (*models.AuditLog, error) {
	log, err := models.NewAuditLog(entityType, action, entityId, "")

	if err != nil {
		return nil, err
	}

	createdLog, err := a.repo.Create(log)
	if err != nil {
		return nil, err
	}
	return createdLog, nil
}

func (a *auditLogService) FindById(id uint) (*models.AuditLog, error) {
	log, err := a.FindById(id)

	if err != nil {
		return nil, err
	}

	return log, nil
}

func (a *auditLogService) GetModelLogs(entityType string, entityId uint) ([]models.AuditLog, error) {
	logs, err := a.repo.GetModelLogs(entityType, entityId)

	if err != nil {
		return nil, err
	}

	return logs, nil
}
