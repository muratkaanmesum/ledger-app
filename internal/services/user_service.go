package services

import (
	"errors"
	"fmt"
	"ptm/internal/models"
	"ptm/internal/repositories"
)

type UserServiceInterface interface {
	RegisterUser(user *models.User) (*models.User, error)
	GetAllUsers(limit, offset int) ([]models.User, error)
	GetUserById(id uint) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
}

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) RegisterUser(user *models.User) (*models.User, error) {
	// Check if the user already exists
	existingUser, err := s.userRepo.GetUserByUsername(user.Username)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists")
	} else if err != nil && err.Error() != "record not found" {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Hash the user's password
	if err := user.HashPassword(); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Save the user to the database
	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetAllUsers(limit, offset int) ([]models.User, error) {
	users, err := s.userRepo.GetUsers(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}
	return users, nil
}

func (s *UserService) GetUserById(id uint) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found with ID %d: %w", id, err)
	}
	return user, nil
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("user not found with username %s: %w", username, err)
	}
	return user, nil
}
