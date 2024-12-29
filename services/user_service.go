package services

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"ptm/db"
	"ptm/db/redis"
	"ptm/models"
	"strconv"
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
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func Send(senderId int, receiverId int, amount int) error {
	var (
		sender   models.User
		receiver models.User
	)

	if err := db.DB.First(&sender, senderId).Error; err != nil {
		return fmt.Errorf("failed to fetch sender: %w", err)
	}
	if err := db.DB.First(&receiver, receiverId).Error; err != nil {
		return fmt.Errorf("failed to fetch receiver: %w", err)
	}

	senderBalanceStr, err := redis.Get(strconv.Itoa(senderId))
	if err != nil {
		return fmt.Errorf("failed to get sender balance from Redis: %w", err)
	}
	receiverBalanceStr, err := redis.Get(strconv.Itoa(receiverId))
	if err != nil {
		return fmt.Errorf("failed to get receiver balance from Redis: %w", err)
	}

	senderBalance, err := strconv.Atoi(senderBalanceStr)
	if err != nil {
		return fmt.Errorf("invalid sender balance in Redis: %w", err)
	}
	receiverBalance, err := strconv.Atoi(receiverBalanceStr)
	if err != nil {
		return fmt.Errorf("invalid receiver balance in Redis: %w", err)
	}

	if senderBalance < amount {
		return errors.New("insufficient balance")
	}

	newSenderBalance := senderBalance - amount
	newReceiverBalance := receiverBalance + amount

	transactionService := NewTransactionService()

	if err := transactionService.CreateTransaction(senderId, -float64(amount), "debit"); err != nil {
		return fmt.Errorf("failed to create sender transaction: %w", err)
	}

	if err := transactionService.CreateTransaction(receiverId, float64(amount), "credit"); err != nil {
		return fmt.Errorf("failed to create receiver transaction: %w", err)
	}

	if err := redis.Set(strconv.Itoa(senderId), strconv.Itoa(newSenderBalance), 0); err != nil {
		return fmt.Errorf("failed to update sender balance in Redis: %w", err)
	}
	if err := redis.Set(strconv.Itoa(receiverId), strconv.Itoa(newReceiverBalance), 0); err != nil {
		return fmt.Errorf("failed to update receiver balance in Redis: %w", err)
	}

	return nil
}
