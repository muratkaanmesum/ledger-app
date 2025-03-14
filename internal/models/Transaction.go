package models

import (
	"errors"
	"time"
)

type Transaction struct {
	ID         uint      `gorm:"primaryKey"`
	FromUserID uint      `gorm:"not null;index"`
	ToUserID   uint      `gorm:"not null;index"`
	Amount     float64   `gorm:"not null"`
	Type       string    `gorm:"not null"`
	Status     string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`

	FromUser User `gorm:"foreignKey:FromUserID;constraint:OnDelete:CASCADE" json:"-"`
	ToUser   User `gorm:"foreignKey:ToUserID;constraint:OnDelete:CASCADE" json:"-"`
}
type TransactionType string

const (
	Credit   TransactionType = "credit"
	Debit    TransactionType = "debit"
	Transfer TransactionType = "transfer"
)

const (
	TransactionStatusPending   = "pending"
	TransactionStatusCompleted = "completed"
	TransactionStatusFailed    = "failed"
)

var validTransactionTypes = map[string]bool{"debit": true, "credit": true}
var validTransactionStatuses = map[string]bool{"pending": true, "completed": true}

func NewTransaction(fromUserID, toUserID uint, amount float64, transactionType TransactionType, status string) (*Transaction, error) {
	if !validTransactionTypes[string(transactionType)] {
		return nil, errors.New("invalid transaction type")
	}
	if !validTransactionStatuses[status] {
		return nil, errors.New("invalid transaction status")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	return &Transaction{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Amount:     amount,
		Type:       string(transactionType),
		Status:     status,
		CreatedAt:  time.Now(),
	}, nil
}
