package repositories

import (
	"errors"
	"gorm.io/gorm"
	"ptm/internal/db"
	"ptm/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
	GetUsers(page, count uint) ([]models.User, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) CreateUser(user *models.User) error {
	return db.DB.Create(user).Error
}

func (r *userRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := db.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(user *models.User) error {
	return db.DB.Save(user).Error
}

func (r *userRepository) DeleteUser(id uint) error {
	return db.DB.Delete(&models.User{}, id).Error
}

func (r *userRepository) GetUsers(page, count uint) ([]models.User, error) {
	var users []models.User
	offset := (page - 1) * count

	err := db.DB.Limit(int(count)).Offset(int(offset)).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
