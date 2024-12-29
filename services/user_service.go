package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"ptm/db"
	"ptm/models"
)

func RegisterUser(user *models.User) (*models.User, error) {
	dbUser := models.User{}
	if err := db.DB.Where("username = ?", user.Username).First(&dbUser).Error; err == nil {
		return nil, errors.New("user already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	if err := db.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserById(id int) (*models.User, error) {
	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := db.DB.Where(&models.User{Username: username}).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with username '%s' not found", username)
		}
		return nil, fmt.Errorf("failed to fetch user with username '%s': %w", username, err)
	}

	return &user, nil
}
