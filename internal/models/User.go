package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"ptm/internal/utils/customError"
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"unique;not null"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *User) VerifyUser(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return customError.Forbidden("The password you provided is incorrect")
	}
	return nil
}

func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.MinCost)
	if err != nil {
		return customError.InternalServerError("Error hashing password")
	}

	u.PasswordHash = string(bytes)
	return nil
}

var ValidRoles = map[string]bool{
	"admin": true,
	"user":  true,
}

func NewUser(username, email, password, role string) (*User, error) {
	if !ValidRoles[role] {
		return nil, errors.New("user role is not valid")
	}
	user := &User{
		Username:     username,
		Email:        email,
		PasswordHash: password,
		Role:         role,
	}
	if err := user.HashPassword(); err != nil {
		return nil, err
	}
	return user, nil
}
