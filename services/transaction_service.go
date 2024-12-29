package services

import (
	"fmt"
	"ptm/db"
	"ptm/models"
)

type TransactionService struct{}

func NewTransactionService() *TransactionService {
	return &TransactionService{}
}

func (t *TransactionService) CreateTransaction(userID int, amount float64, transactionType string) error {
	transaction := models.Transaction{
		ID:     uint(userID),
		Amount: amount,
		Type:   transactionType,
	}

	if err := db.DB.Create(&transaction).Error; err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
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
