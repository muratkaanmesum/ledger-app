package services

import (
	"ptm/db"
	"ptm/models"
)

func CreateTransaction(userID uint, amount float64, txType string) (*models.Transaction, error) {
	transaction := &models.Transaction{UserID: userID, Amount: amount, Type: txType}
	if err := db.DB.Create(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func ListTransactions(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := db.DB.Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
