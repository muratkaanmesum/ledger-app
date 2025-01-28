package services

import (
	"ptm/internal/models"
	"ptm/internal/repositories"
)

type TransactionService interface {
	CreateTransaction(fromId, toId uint, amount float64, transactionType models.TransactionType) (*models.Transaction, error)
	ListTransactions(userID uint, page, count uint, failed bool) ([]models.Transaction, error)
	UpdateTransactionState(transactionId uint, state models.TransactionType) error
	GetTransactionById(transactionId uint) (*models.Transaction, error)
}

type transactionService struct {
	repository repositories.TransactionRepository
}

func NewTransactionService(repository repositories.TransactionRepository) TransactionService {
	return &transactionService{
		repository: repository,
	}
}

func (t *transactionService) CreateTransaction(fromId, toId uint, amount float64, transactionType models.TransactionType) (*models.Transaction, error) {
	toTransaction, err := models.NewTransaction(fromId, toId, amount, transactionType, models.TransactionStatusPending)
	if err != nil {
		return nil, err
	}

	if err := t.repository.CreateTransaction(toTransaction); err != nil {
		return nil, err
	}

	return toTransaction, nil
}

func (t *transactionService) UpdateTransactionState(transactionId uint, state models.TransactionType) error {
	transaction, err := t.repository.GetTransactionByID(transactionId)

	if err != nil {
		return err
	}

	transaction.Status = string(state)
	if err := t.repository.UpdateTransaction(transaction); err != nil {
		return err
	}

	return nil
}

func (t *transactionService) ListTransactions(userID uint, page, count uint, failed bool) ([]models.Transaction, error) {
	var transactions []models.Transaction

	transactions, err := t.repository.GetAllTransactions(userID, page, count, failed)

	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (t *transactionService) GetTransactionById(transactionId uint) (*models.Transaction, error) {
	transaction, err := t.repository.GetTransactionByID(transactionId)

	if err != nil {
		return nil, err
	}
	return transaction, nil
}
