package services

import (
	"errors"
	"ptm/internal/models"
	"ptm/internal/repositories"
)

type BalanceService interface {
	GetUserBalance(userID uint) (*models.Balance, error)
	UpdateUserBalance(userID uint, amount float64) error
	IncrementUserBalance(userID uint, amount float64) error
	DecrementUserBalance(userID uint, amount float64) error
	CreateBalance(user *models.User) (*models.Balance, error)
}
type balanceService struct {
	repo              repositories.BalanceRepository
	historyRepository repositories.BalanceHistoryRepository
}

func NewBalanceService(balanceRepository repositories.BalanceRepository, historyRepository repositories.BalanceHistoryRepository) BalanceService {
	return &balanceService{
		repo: balanceRepository,
	}
}

func (s *balanceService) CreateBalance(user *models.User) (*models.Balance, error) {
	exists, err := s.repo.GetBalance(user.ID)
	if exists != nil {
		return nil, err
	}

	balance, createErr := s.repo.CreateBalance(user.ID, 0)
	if createErr != nil {
		return nil, createErr
	}

	return balance, nil
}

func (s *balanceService) GetUserBalance(userID uint) (*models.Balance, error) {
	balance, err := s.repo.GetBalance(userID)
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
	if err := s.repo.UpdateBalance(userID, amount); err != nil {
		return err
	}
	if err := s.historyRepository.Create(models.NewBalanceHistory(userID, amount)); err != nil {
		return err
	}

	return nil
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

func (s *balanceService) GetBalanceAtTime(userID uint) (*models.Balance, error) {
	balance, err := s.repo.GetBalance(userID)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
