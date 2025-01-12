package transaction

import (
	"gorm.io/gorm"
	"ptm/internal/db"
)

func StartTransaction() (*gorm.DB, error) {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func CommitTransaction(tx *gorm.DB) error {
	if tx == nil {
		return nil // No transaction to commit
	}
	return tx.Commit().Error
}

func RollbackTransaction(tx *gorm.DB) error {
	if tx == nil {
		return nil // No transaction to rollback
	}
	return tx.Rollback().Error
}
