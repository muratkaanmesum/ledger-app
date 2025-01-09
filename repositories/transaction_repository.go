package repositories

import (
	"gorm.io/gorm"
	"ptm/db" // Assuming this package contains your global DB instance
	"ptm/models"
)

type TransactionRepository struct {
	tx *gorm.DB
}

type ITransactionRepository interface {
	StartTransaction() error
	CommitTransaction() error
	RollbackTransaction() error
	CreateTransaction(transaction *models.Transaction) error
	GetTransactionByID(id uint) (*models.Transaction, error)
	GetAllTransactions() ([]models.Transaction, error)
	UpdateTransaction(transaction *models.Transaction) error
	DeleteTransaction(id uint) error
}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}

func (r *TransactionRepository) StartTransaction() error {
	if r.tx != nil {
		return nil
	}
	r.tx = db.DB.Begin()
	return r.tx.Error
}

func (r *TransactionRepository) CommitTransaction() error {
	if r.tx == nil {
		return nil
	}
	err := r.tx.Commit().Error
	r.tx = nil
	return err
}

func (r *TransactionRepository) RollbackTransaction() error {
	if r.tx == nil {
		return nil
	}
	err := r.tx.Rollback().Error
	r.tx = nil
	return err
}

func (r *TransactionRepository) CreateTransaction(transaction *models.Transaction) error {
	if r.tx != nil {
		return r.tx.Create(transaction).Error
	}
	return db.DB.Create(transaction).Error
}

func (r *TransactionRepository) GetTransactionByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	var err error
	if r.tx != nil {
		err = r.tx.First(&transaction, id).Error
	} else {
		err = db.DB.First(&transaction, id).Error
	}
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) GetAllTransactions() ([]models.Transaction, error) {
	var transactions []models.Transaction
	var err error
	if r.tx != nil {
		err = r.tx.Find(&transactions).Error
	} else {
		err = db.DB.Find(&transactions).Error
	}
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) UpdateTransaction(transaction *models.Transaction) error {
	if r.tx != nil {
		return r.tx.Save(transaction).Error
	}
	return db.DB.Save(transaction).Error
}

func (r *TransactionRepository) DeleteTransaction(id uint) error {
	if r.tx != nil {
		return r.tx.Delete(&models.Transaction{}, id).Error
	}
	return db.DB.Delete(&models.Transaction{}, id).Error
}
