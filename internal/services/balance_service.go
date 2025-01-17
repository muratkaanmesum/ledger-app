package services

import (
	"errors"
	"ptm/internal/models"
	"ptm/internal/repositories"
	"time"
)

type BalanceService interface {
	GetUserBalance(userID uint) (*models.Balance, error)
	UpdateUserBalance(userID uint, amount float64) error
	IncrementUserBalance(userID uint, amount float64) error
	DecrementUserBalance(userID uint, amount float64) error
	CreateBalance(user *models.User) (*models.Balance, error)
}
type balanceService struct {
	repo repositories.BalanceRepository
}

func NewBalanceService(repo repositories.BalanceRepository) BalanceService {
	return &balanceService{repo: repo}
}

func (s *balanceService) CreateBalance(user *models.User) (*models.Balance, error) {
	balance, err := s.repo.CreateBalance(user.ID, 0)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (s *balanceService) GetUserBalance(userID uint) (*models.Balance, error) {
	balance, err := s.repo.GetBalance(userID, nil)
	if err != nil {
		return nil, err
	}
	if balance == nil {
		return nil, errors.New("balance not found")
	}
	return balance, nil
}

func (s *balanceService) UpdateUserBalance(userID uint, amount float64) error {
	if amount < 0 {
		return errors.New("amount must be non-negative")
	}
	return s.repo.UpdateBalance(userID, amount)
}

func (s *balanceService) IncrementUserBalance(userID uint, amount float64) error {
	if amount <= 0 {
		return errors.New("increment amount must be greater than zero")
	}
	return s.repo.IncrementBalance(userID, amount)
}

func (s *balanceService) DecrementUserBalance(userID uint, amount float64) error {
	if amount <= 0 {
		return errors.New("decrement amount must be greater than zero")
	}
	return s.repo.DecrementBalance(userID, amount)
}

func (s *balanceService) GetBalanceAtTime(userID uint, date *time.Time) (*models.Balance, error) {
	balance, err := s.repo.GetBalance(userID, date)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
