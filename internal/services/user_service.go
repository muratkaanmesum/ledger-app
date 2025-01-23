package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"ptm/internal/models"
	"ptm/internal/repositories"
)

type UserService interface {
	RegisterUser(user *models.User) (*models.User, error)
	GetAllUsers(limit, offset int) ([]models.User, error)
	GetUserById(id uint) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	Exists(userId uint) (bool, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {

	return &userService{
		userRepo: userRepository,
	}
}

func (s *userService) RegisterUser(user *models.User) (*models.User, error) {
	existingUser, err := s.userRepo.GetUserByUsername(user.Username)
	fmt.Println(existingUser, "TEST", err)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}
	if err := user.HashPassword(); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *userService) GetAllUsers(limit, offset int) ([]models.User, error) {
	users, err := s.userRepo.GetUsers(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}
	return users, nil
}

func (s *userService) GetUserById(id uint) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found with ID %d: %w", id, err)
	}
	return user, nil
}

func (s *userService) GetUserByUsername(username string) (*models.User, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("user not found with username %s: %w", username, err)
	}
	return user, nil
}

func (s *userService) Exists(userId uint) (bool, error) {
	_, err := s.userRepo.GetUserByID(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}
	return true, nil
}
