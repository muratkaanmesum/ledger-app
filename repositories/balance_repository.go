package repositories

import (
	"errors"
	"sync"
	"time"

	"gorm.io/gorm"
	"ptm/db"
	"ptm/models"
)

type BalanceRepository interface {
	GetBalance(userID uint, date *time.Time) (*models.Balance, error)
	UpdateBalance(userID uint, amount float64) error
	IncrementBalance(userID uint, amount float64) error
	DecrementBalance(userID uint, amount float64) error
}

type balanceRepository struct {
	mu sync.RWMutex
}

func (r *balanceRepository) GetBalance(userID uint, date *time.Time) (*models.Balance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var balance models.Balance
	query := db.DB.Where("user_id = ?", userID)

	if date != nil {
		query = query.Where("last_updated_at <= ?", date).Order("last_updated_at DESC")
	} else {
		query = query.Order("last_updated_at DESC")
	}

	err := query.First(&balance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &balance, nil
}

func (r *balanceRepository) UpdateBalance(userID uint, amount float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var balance models.Balance
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).First(&balance).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				balance = models.Balance{
					UserID:        userID,
					Amount:        amount,
					LastUpdatedAt: time.Now(),
				}
				return tx.Create(&balance).Error
			}
			return err
		}

		balance.Amount = amount
		balance.LastUpdatedAt = time.Now()
		return tx.Save(&balance).Error
	})

	return err
}

func (r *balanceRepository) IncrementBalance(userID uint, amount float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var balance models.Balance
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).First(&balance).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				balance = models.Balance{
					UserID:        userID,
					Amount:        amount,
					LastUpdatedAt: time.Now(),
				}
				return tx.Create(&balance).Error
			}
			return err
		}

		balance.Amount += amount
		balance.LastUpdatedAt = time.Now()
		return tx.Save(&balance).Error
	})

	return err
}

func (r *balanceRepository) DecrementBalance(userID uint, amount float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var balance models.Balance
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).First(&balance).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("balance not found")
			}
			return err
		}

		if balance.Amount < amount {
			return errors.New("insufficient balance")
		}

		balance.Amount -= amount
		balance.LastUpdatedAt = time.Now()
		return tx.Save(&balance).Error
	})

	return err
}
