package services

import (
	"ptm/db"
	"ptm/models"
)

type TransactionService struct{}

func NewTransactionService() *TransactionService {
	return &TransactionService{}
}

func (t *TransactionService) CreateTransaction(fromId, toId uint, amount float64, transactionType string) error {
	transaction, err := models.NewTransaction(fromId, toId, amount, transactionType, models.TransactionStatusPending)
	if err != nil {
		return err
	}

	if err := db.DB.Create(&transaction).Error; err != nil {
		return err
	}
	return nil
}

func ListTransactions(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := db.DB.Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
