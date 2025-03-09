package repositories

import (
	"errors"
	"ptm/internal/db"
	"ptm/internal/models"
	"ptm/pkg/utils/customError"
	"sync"
	"time"

	"gorm.io/gorm"
)

var (
	instance *balanceRepository
	once     sync.Once
)

type BalanceRepository interface {
	GetBalance(userID uint) (*models.Balance, error)
	UpdateBalance(userID uint, amount float64) (*models.Balance, error)
	IncrementBalance(userID uint, amount float64) (*models.Balance, error)
	DecrementBalance(userID uint, amount float64) (*models.Balance, error)
	CreateBalance(userID uint, amount float64) (*models.Balance, error)
}

type balanceRepository struct {
	mu *sync.RWMutex
}

func NewBalanceRepository() BalanceRepository {
	once.Do(func() {
		instance = &balanceRepository{
			mu: &sync.RWMutex{},
		}
	})
	return instance
}

func (r *balanceRepository) GetBalance(userID uint) (*models.Balance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var balance models.Balance
	query := db.DB.Where("user_id = ?", userID).First(&balance)

	if query.Error != nil {
		return nil, customError.NotFound("Balance not found", query.Error)
	}

	return &balance, nil
}

func (r *balanceRepository) UpdateBalance(userID uint, amount float64) (*models.Balance, error) {
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

	return &balance, err
}

func (r *balanceRepository) IncrementBalance(userID uint, amount float64) (*models.Balance, error) {
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

	return &balance, err
}

func (r *balanceRepository) DecrementBalance(userID uint, amount float64) (*models.Balance, error) {
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

	return &balance, err
}

func (r *balanceRepository) CreateBalance(userID uint, amount float64) (*models.Balance, error) {
	var balance models.Balance
	balanceDB := db.DB.Create(&models.Balance{
		UserID: userID,
		Amount: amount,
	})
	if balanceDB.Error != nil {
		return nil, balanceDB.Error
	}

	return &balance, nil
}
