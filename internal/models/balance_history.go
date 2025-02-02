package models

import "time"

type BalanceHistory struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"index"`
	Amount     float64   `gorm:"type:decimal(10,2)"`
	ChangeType string    `gorm:"type:varchar(50)"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

func NewBalanceHistory(userID uint, amount float64) *BalanceHistory {
	changeType := ""
	if amount > 0 {
		changeType = "credit"
	} else {
		changeType = "debit"
	}
	return &BalanceHistory{
		UserID:     userID,
		Amount:     amount,
		ChangeType: changeType,
		CreatedAt:  time.Now(),
	}
}
