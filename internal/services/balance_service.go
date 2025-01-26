package services

import (
	"errors"
	"ptm/internal/models"
	"ptm/internal/repositories"
	"ptm/internal/utils/customError"
	"time"
)

type BalanceService interface {
	GetUserBalance(userID uint) (*models.Balance, error)
	UpdateUserBalance(userID uint, amount float64) error
	IncrementUserBalance(userID uint, amount float64) error
	DecrementUserBalance(userID uint, amount float64) error
	CreateBalance(user *models.User) (*models.Balance, error)
	GetBalanceAtTime(userID uint, time time.Time) (*models.BalanceHistory, error)
}
type balanceService struct {
	repo              repositories.BalanceRepository
	historyRepository repositories.BalanceHistoryRepository
}

func NewBalanceService(balanceRepository repositories.BalanceRepository, historyRepository repositories.BalanceHistoryRepository) BalanceService {
	return &balanceService{
		repo:              balanceRepository,
		historyRepository: historyRepository,
	}
}

func (s *balanceService) CreateBalance(user *models.User) (*models.Balance, error) {
	exists, err := s.repo.GetBalance(user.ID)
	if exists != nil {
		return nil, customError.NotFound("User not found", err)
	}

	balance, createErr := s.repo.CreateBalance(user.ID, 0)
	if createErr != nil {
		return nil, customError.InternalServerError("Failed to create balance", createErr)
	}

	return balance, nil
}

func (s *balanceService) GetUserBalance(userID uint) (*models.Balance, error) {
	balance, err := s.repo.GetBalance(userID)
	if err != nil {
		return nil, err
	}
	if balance == nil {
		return nil, customError.NotFound("User not found")
	}
	return balance, nil
}

func (s *balanceService) UpdateUserBalance(userID uint, amount float64) error {
	if amount < 0 {
		return customError.BadRequest("Amount must be greater than zero")
	}
	if err := s.repo.UpdateBalance(userID, amount); err != nil {
		return customError.InternalServerError("Failed to update balance", err)
	}
	if err := s.historyRepository.Create(models.NewBalanceHistory(userID, amount)); err != nil {
		return customError.InternalServerError("Failed to Create history for balance", err)
	}

	return nil
}

func (s *balanceService) IncrementUserBalance(userID uint, amount float64) error {
	if amount <= 0 {
		return errors.New("increment amount must be greater than zero")
	}

	if err := s.repo.IncrementBalance(userID, amount); err != nil {
		return customError.InternalServerError("Failed to increment balance", err)
	}

	if err := s.historyRepository.Create(models.NewBalanceHistory(userID, amount)); err != nil {
		return customError.InternalServerError("Failed to Create history for balance", err)
	}

	return nil
}
func (s *balanceService) DecrementUserBalance(userID uint, amount float64) error {
	if amount <= 0 {
		return customError.BadRequest("Amount must be greater than zero")
	}

	if err := s.repo.DecrementBalance(userID, amount); err != nil {
		return customError.InternalServerError("Failed to Decrement balance", err)
	}

	if err := s.historyRepository.Create(models.NewBalanceHistory(userID, amount)); err != nil {
		return customError.InternalServerError("Failed to Create history for balance", err)
	}

	return nil
}

func (s *balanceService) GetBalanceAtTime(userID uint, time time.Time) (*models.BalanceHistory, error) {
	history, err := s.historyRepository.GetBalanceAtTime(userID, time)
	if err != nil {
		return nil, customError.InternalServerError("Failed to Get balance at time", err)
	}

	return history, nil
}
