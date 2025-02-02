package repositories

import (
	"ptm/internal/db"
	"ptm/internal/models"
	"time"
)

type balanceHistoryRepository struct{}

type BalanceHistoryRepository interface {
	Create(entry *models.BalanceHistory) error
	GetByID(id uint) (*models.BalanceHistory, error)
	GetByUserID(userID string, limit, offset int) ([]models.BalanceHistory, error)
	Update(entry *models.BalanceHistory) error
	Delete(id uint) error
	GetBalanceAtTime(userId uint, time time.Time) (*models.BalanceHistory, error)
	GetUserHistories(userId uint) ([]models.BalanceHistory, error)
}

func NewBalanceHistoryRepository() BalanceHistoryRepository {
	return &balanceHistoryRepository{}
}

func (r *balanceHistoryRepository) Create(entry *models.BalanceHistory) error {
	return db.DB.Create(entry).Error
}

func (r *balanceHistoryRepository) GetByID(id uint) (*models.BalanceHistory, error) {
	var history models.BalanceHistory
	err := db.DB.First(&history, id).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

func (r *balanceHistoryRepository) GetByUserID(userID string, limit, offset int) ([]models.BalanceHistory, error) {
	var histories []models.BalanceHistory
	err := db.DB.Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

func (r *balanceHistoryRepository) Update(entry *models.BalanceHistory) error {
	return db.DB.Save(entry).Error
}

func (r *balanceHistoryRepository) Delete(id uint) error {
	return db.DB.Delete(&models.BalanceHistory{}, id).Error
}

func (r *balanceHistoryRepository) GetBalanceAtTime(userId uint, time time.Time) (*models.BalanceHistory, error) {
	history := &models.BalanceHistory{}

	err := db.DB.Where("created_at < ? and user_id = ?", time, userId).
		Find(&models.BalanceHistory{}).First(&history).Error
	if err != nil {
		return nil, err
	}

	return history, nil
}

func (r *balanceHistoryRepository) GetUserHistories(userId uint) ([]models.BalanceHistory, error) {
	var history []models.BalanceHistory

	err := db.DB.Where("user_id = ?", userId).
		Order("created_at desc").
		Find(&history).Error

	if err != nil {
		return nil, err
	}

	return history, nil
}
