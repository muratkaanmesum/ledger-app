package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string        `json:"name"`
	Transactions []Transaction `json:"transactions" gorm:"foreignKey:UserID"`
	Role         string        `json:"role"`
}
