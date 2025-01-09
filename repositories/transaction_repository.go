package repositories

import (
	"gorm.io/gorm"
	"ptm/db" // Assuming this package contains your global DB instance
	"ptm/models"
)

type TransactionRepository struct{}

type ITransactionRepository interface {
	CreateTransaction(tx *gorm.DB, transaction *models.Transaction) error
	GetTransactionByID(id uint) (*models.Transaction, error)
	GetAllTransactions() ([]models.Transaction, error)
	UpdateTransaction(tx *gorm.DB, transaction *models.Transaction) error
	DeleteTransaction(tx *gorm.DB, id uint) error
}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}

func (r *TransactionRepository) CreateTransaction(transaction *models.Transaction) error {
	if db.DB != nil {
		return db.DB.Create(transaction).Error
	}
	return db.DB.Create(transaction).Error
}

func (r *TransactionRepository) GetTransactionByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := db.DB.First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) GetAllTransactions() ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := db.DB.Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) UpdateTransaction(transaction *models.Transaction) error {
	if db.DB != nil {
		return db.DB.Save(transaction).Error
	}
	return db.DB.Save(transaction).Error
}

func (r *TransactionRepository) DeleteTransaction(id uint) error {
	if db.DB != nil {
		return db.DB.Delete(&models.Transaction{}, id).Error
	}
	return db.DB.Delete(&models.Transaction{}, id).Error
}
