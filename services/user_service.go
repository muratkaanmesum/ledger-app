package services

import (
	"errors"
	"fmt"
	"ptm/db"

	"gorm.io/gorm"
	"ptm/models"
)

type UserServiceInterface interface {
	RegisterUser(user *models.User) (*models.User, error)
	GetAllUsers(limit, offset int) ([]models.User, error)
	GetUserById(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
}

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) RegisterUser(user *models.User) (*models.User, error) {
	existingUser := models.User{}
	if err := db.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if err := user.HashPassword(); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	if err := db.DB.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetAllUsers(limit, offset int) ([]models.User, error) {
	var users []models.User
	if err := db.DB.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}
	return users, nil
}

func (s *UserService) GetUserById(id int) (*models.User, error) {
	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("user not found with ID %d: %w", id, err)
	}
	return &user, nil
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found with username %s: %w", username, err)
	}
	return &user, nil
}
