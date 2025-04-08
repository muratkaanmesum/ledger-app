package services

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"ptm/internal/db/redis"
	"ptm/internal/models"
	"ptm/internal/repositories"
	"ptm/pkg/logger"
	"ptm/pkg/utils/customError"
	"strconv" // Added import
	"time"
)

type BalanceService interface {
	GetUserBalance(userID uint) (*models.Balance, error)
	UpdateUserBalance(userID uint, amount float64) error
	IncrementUserBalance(userID uint, amount float64) error
	DecrementUserBalance(userID uint, amount float64) error
	CreateBalance(user *models.User) (*models.Balance, error)
	GetBalanceAtTime(userID uint, time time.Time) (*models.BalanceHistory, error)
	GetUserBalanceHistory(userID uint) ([]models.BalanceHistory, error)
}
type balanceService struct {
	repo              repositories.BalanceRepository
	historyRepository repositories.BalanceHistoryRepository
	logService        AuditLogService
}

func NewBalanceService(
	balanceRepository repositories.BalanceRepository,
	historyRepository repositories.BalanceHistoryRepository,
	logService AuditLogService,
) BalanceService {
	return &balanceService{
		repo:              balanceRepository,
		historyRepository: historyRepository,
		logService:        logService,
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
	_, err = s.logService.CreateLog("balance", "create", balance.Id)
	if err != nil {
		logger.Logger.Error("Couldn't log the creation", zap.Error(err))
	}
	return balance, nil
}

func (s *balanceService) GetUserBalance(userID uint) (*models.Balance, error) {
	key := redis.Key("balance", userID)
	exists, err := redis.Exists(key)
	if err != nil {
		return nil, customError.InternalServerError("Internal server error", err)
	}

	if !exists {
		logger.Logger.Info("Balance not exists on redis for", zap.String("user_id", fmt.Sprint(userID)))
		balance, err := s.repo.GetBalance(userID)

		if err != nil {
			return nil, err
		}
		return balance, nil
	}

	balance, err := redis.Get(key)
	if err != nil {
		return nil, customError.InternalServerError("Internal server error", err)
	}
	return models.NewBalance(userID, balance), nil // Updated line
}

func (s *balanceService) UpdateUserBalance(userID uint, amount float64) error {
	if amount < 0 {
		return customError.BadRequest("Amount must be greater than zero")
	}
	balance, err := s.repo.UpdateBalance(userID, amount)

	if err != nil {
		return customError.InternalServerError("Failed to update balance", err)
	}

	key := redis.Key("balance", userID)
	if err := redis.Set(key, balance.Amount); err != nil {
		return customError.InternalServerError("Failed to update balance", err)
	}

	if err := s.historyRepository.Create(models.NewBalanceHistory(userID, balance.Amount)); err != nil {
		return customError.InternalServerError("Failed to Create history for balance", err)
	}

	_, err = s.logService.CreateLog("balance", "create", balance.Id)
	if err != nil {
		logger.Logger.Error("Couldn't log the creation", zap.Error(err))
	}

	return nil
}

func (s *balanceService) IncrementUserBalance(userID uint, amount float64) error {
	if amount <= 0 {
		return errors.New("increment amount must be greater than zero")
	}

	balance, err := s.repo.IncrementBalance(userID, amount)

	if err != nil {
		return customError.InternalServerError("Failed to increment balance", err)
	}

	key := redis.Key("balance", userID)
	if err := redis.Set(key, balance.Amount); err != nil {
		return customError.InternalServerError("Failed to update balance", err)
	}

	if err := s.historyRepository.Create(models.NewBalanceHistory(userID, balance.Amount)); err != nil {
		return customError.InternalServerError("Failed to Create history for balance", err)
	}

	_, err = s.logService.CreateLog("balance", "update", balance.Id)

	if err != nil {
		logger.Logger.Error("Couldn't log the update", zap.Error(err))
	}

	return nil
}
func (s *balanceService) DecrementUserBalance(userID uint, amount float64) error {
	if amount <= 0 {
		return customError.BadRequest("Amount must be greater than zero")
	}

	balance, err := s.repo.DecrementBalance(userID, amount)

	if err != nil {
		return customError.InternalServerError("Failed to Decrement balance", err)
	}

	key := redis.Key("balance", userID)
	if err := redis.Set(key, balance.Amount); err != nil {
		return customError.InternalServerError("Failed to update balance", err)
	}

	if err := s.historyRepository.Create(models.NewBalanceHistory(userID, balance.Amount)); err != nil {
		return customError.InternalServerError("Failed to Create history for balance", err)
	}

	_, err = s.logService.CreateLog("balance", "create", balance.Id)
	if err != nil {
		logger.Logger.Error("Couldn't log the creation", zap.Error(err))
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

func (s *balanceService) GetUserBalanceHistory(userID uint) ([]models.BalanceHistory, error) {

	histories, err := s.historyRepository.GetUserHistories(userID)

	if err != nil {
		return nil, err
	}

	return histories, nil
}

func parseAmount(amount string) float64 { // Added helper function
	parsed, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0.0
	}
	return parsed
}
