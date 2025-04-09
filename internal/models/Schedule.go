package models

import "time"

type Schedule struct {
	ID           uint    `gorm:"primaryKey"`
	UserID       uint    `gorm:"not null"`
	TargetUserID uint    `gorm:"not null"`
	Amount       float64 `gorm:"not null"`
	//Currency     string    `gorm:"size:10;not null"`
	Status          string `gorm:"size:20;default:'pending'"`
	TransactionType TransactionType
	ExecuteAt       time.Time `gorm:"not null"`
	User            User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}
