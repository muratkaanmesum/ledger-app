package models

import (
	"strconv"
	"time"
)

type Balance struct {
	Id            uint    `gorm:"primaryKey"`
	UserID        uint    `gorm:"not null"`
	Amount        float64 `gorm:"default:0.0"`
	LastUpdatedAt time.Time
	User          User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

func NewBalance(userID uint, amount string) *Balance {
	return &Balance{
		Id:     0,
		UserID: userID,
		Amount: parseAmount(amount),
	}
}

func parseAmount(amount string) float64 {
	parsed, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0.0
	}
	return parsed
}
