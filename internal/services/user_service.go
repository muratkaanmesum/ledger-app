package services

import (
	"errors"
	"gorm.io/gorm"
	"ptm/internal/dtos"
	"ptm/internal/models"
	"ptm/internal/repositories"
	"ptm/pkg/utils/customError"
)

type UserService interface {
	RegisterUser(user *models.User) (*models.User, error)
	GetAllUsers(page, count uint) ([]models.User, error)
	GetUserById(id uint) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	Exists(userId uint) (bool, error)
	UpdateUser(id uint, user *dtos.UpdateUserRequest) (*models.User, error)
	DeleteUser(userId uint) error
	GetUserRules(userID uint) (models.Rule, error) // New method added
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
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, customError.BadRequest("User already exists")
	}
	if err := user.HashPassword(); err != nil {
		return nil, customError.InternalServerError("failed to hash password")
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, customError.InternalServerError("failed to create user")
	}

	return user, nil
}

func (s *userService) GetAllUsers(page, count uint) ([]models.User, error) {
	users, err := s.userRepo.GetUsers(page, count)
	if err != nil {
		return nil, customError.InternalServerError("failed to get users")
	}
	return users, nil
}

func (s *userService) GetUserById(id uint) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, customError.NotFound("User not found")
	}
	return user, nil
}

func (s *userService) GetUserByUsername(username string) (*models.User, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Exists(userId uint) (bool, error) {
	_, err := s.userRepo.GetUserByID(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, customError.InternalServerError("failed to get user")
	}
	return true, nil
}

func (s *userService) DeleteUser(userId uint) error {
	if err := s.userRepo.DeleteUser(userId); err != nil {
		return err
	}

	return nil
}

func (s *userService) UpdateUser(id uint, user *dtos.UpdateUserRequest) (*models.User, error) {
	existingUser, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	existingUser.Username = user.Username
	existingUser.Email = user.Email

	existingUser.PasswordHash = user.Password

	if err := existingUser.HashPassword(); err != nil {
		return nil, customError.InternalServerError("failed to hash password")
	}

	if err := s.userRepo.UpdateUser(existingUser); err != nil {
		return nil, customError.InternalServerError("Failed to update user")
	}

	return existingUser, nil
}

func (s *userService) GetUserRules(userID uint) (models.Rule, error) {
	return s.userRepo.GetUserRules(userID)
}
