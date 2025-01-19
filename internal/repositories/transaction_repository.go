package repositories

import (
	"ptm/internal/db"
	"ptm/internal/models"
)

type transactionRepository struct{}

type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) error
	GetTransactionByID(id uint) (*models.Transaction, error)
	GetAllTransactions() ([]models.Transaction, error)
	UpdateTransaction(transaction *models.Transaction) error
	DeleteTransaction(id uint) error
}

func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{}
}

func (r *transactionRepository) CreateTransaction(transaction *models.Transaction) error {
	return db.DB.Create(transaction).Error
}

func (r *transactionRepository) GetTransactionByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := db.DB.First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) GetAllTransactions() ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := db.DB.Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) UpdateTransaction(transaction *models.Transaction) error {
	if db.DB != nil {
		return db.DB.Save(transaction).Error
	}
	return db.DB.Save(transaction).Error
}

func (r *transactionRepository) DeleteTransaction(id uint) error {
	if db.DB != nil {
		return db.DB.Delete(&models.Transaction{}, id).Error
	}
	return db.DB.Delete(&models.Transaction{}, id).Error
}
