package models

import (
	"strconv"
	"time"
)

type Balance struct {
	Id            uint    `gorm:"primaryKey"`
	UserID        uint    `gorm:"not null"`
	Amount        float64 `gorm:"default:0.0"`
	Currency      string  `gorm:"size:3;not null;default:TRY"`
	LastUpdatedAt time.Time
	User          User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

func NewBalanceFromString(userID uint, amount string) *Balance {
	return &Balance{
		Id:     0,
		UserID: userID,
		Amount: parseAmount(amount),
	}
}

func NewBalanceFromFloat(userID uint, amount float64) *Balance {
	return &Balance{
		Id:     0,
		UserID: userID,
		Amount: amount,
	}
}

func parseAmount(amount string) float64 {
	parsed, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0.0
	}
	return parsed
}
