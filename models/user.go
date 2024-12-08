package models

import "gorm.io/gorm"

// User model definition
type User struct {
	gorm.Model        // Adds ID, CreatedAt, UpdatedAt, DeletedAt fields
	Name       string `json:"name"`
	Email      string `json:"email" gorm:"unique"`
}
