package services

import (
	"ptm/internal/db"
	"ptm/internal/models"
	"ptm/internal/repositories"
)

type TransactionServiceInterface interface {
	CreateTransaction(fromId, toId uint, amount float64, transactionType models.TransactionType) (*models.Transaction, error)
	ListTransactions(userID uint) ([]models.Transaction, error)
	UpdateTransactionState(fromId uint) error
}

type TransactionService struct {
	repository repositories.TransactionRepository
}

func NewTransactionService() TransactionServiceInterface {
	return &TransactionService{}
}

func (t *TransactionService) CreateTransaction(fromId, toId uint, amount float64, transactionType models.TransactionType) (*models.Transaction, error) {
	toTransaction, err := models.NewTransaction(fromId, toId, amount, transactionType, models.TransactionStatusPending)
	if err != nil {
		return nil, err
	}

	if err := t.repository.CreateTransaction(toTransaction); err != nil {
		return nil, err
	}

	return toTransaction, nil
}

func (t *TransactionService) UpdateTransactionState(fromId uint) error {
	transaction, err := t.repository.GetTransactionByID(fromId)

	if err != nil {
		return err
	}

	transaction.Status = models.TransactionStatusPending
	if err := t.repository.UpdateTransaction(transaction); err != nil {
		return err
	}

	return nil
}

func (t *TransactionService) ListTransactions(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := db.DB.Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
