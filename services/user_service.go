package services

import (
	"errors"
	"fmt"
	"ptm/db"
	"ptm/db/redis"
	"ptm/models"
	"strconv"
)

func CreateUser(name, role string) (*models.User, error) {
	user := &models.User{Name: name, Role: role}
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

	senderTransaction := models.Transaction{
		UserID: uint(senderId),
		Amount: -float64(amount),
		Type:   "debit",
	}
	receiverTransaction := models.Transaction{
		UserID: uint(receiverId),
		Amount: float64(amount),
		Type:   "credit",
	}

	if err := db.DB.Create(&senderTransaction).Error; err != nil {
		return fmt.Errorf("failed to create sender transaction: %w", err)
	}
	if err := db.DB.Create(&receiverTransaction).Error; err != nil {
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
