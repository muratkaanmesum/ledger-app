package models

import (
	"golang.org/x/crypto/bcrypt"
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
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err == nil {
		return err
	}

	return nil
}

func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.MinCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(bytes)
	return nil
}

func (u *User) Validate() {

}
