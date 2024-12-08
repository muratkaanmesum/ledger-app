package db

import (
	"log"
	"math/rand"
	"ptm/db/redis"
	"ptm/models"
	"ptm/services"
	"strconv"
)

func SeedUsers() {
	users := []models.User{
		{Name: "John Doe", Role: "admin"},
		{Name: "Jane Smith", Role: "normal"},
		{Name: "Alice Johnson", Role: "normal"},
	}

	transactionService := services.NewTransactionService()

	for _, user := range users {
		var existingUser models.User
		if err := DB.Where("name = ?", user.Name).First(&existingUser).Error; err == nil {
			log.Printf("User with name %s already exists. Skipping seed.", user.Name)
			continue
		}

		if err := DB.Create(&user).Error; err != nil {
			log.Printf("Failed to seed user %s: %v", user.Name, err)
			continue
		}

		userID := int(user.ID)
		initialBalance := rand.Intn(1000) + 100

		if err := redis.Set(strconv.Itoa(userID), strconv.Itoa(initialBalance), 0); err != nil {
			log.Printf("Failed to set Redis balance for user %s: %v", user.Name, err)
		} else {
			log.Printf("Initialized Redis balance for user %s with balance %d", user.Name, initialBalance)
		}

		if err := transactionService.CreateTransaction(userID, float64(initialBalance), "deposit"); err != nil {
			log.Printf("Failed to create initial deposit transaction for user %s: %v", user.Name, err)
		} else {
			log.Printf("Created initial deposit transaction for user %s", user.Name)
		}
	}

	createSpecificTransactions(users, transactionService)
}

func createSpecificTransactions(users []models.User, transactionService *services.TransactionService) {
	userCount := len(users)
	if userCount < 2 {
		log.Println("Not enough users to create send/receive transactions")
		return
	}

	for _, sender := range users {
		senderID := int(sender.ID)

		depositAmount := float64(rand.Intn(500) + 50)
		if err := transactionService.CreateTransaction(senderID, depositAmount, "deposit"); err != nil {
			log.Printf("Failed to create deposit transaction for user %d: %v", senderID, err)
		} else {
			log.Printf("Created deposit transaction of %.2f for user %d", depositAmount, senderID)
		}

		withdrawAmount := float64(rand.Intn(200) + 10)
		if err := transactionService.CreateTransaction(senderID, -withdrawAmount, "withdraw"); err != nil {
			log.Printf("Failed to create withdraw transaction for user %d: %v", senderID, err)
		} else {
			log.Printf("Created withdraw transaction of %.2f for user %d", withdrawAmount, senderID)
		}

		amount := float64(rand.Intn(200) + 20)
		receiver := selectRandomUser(users, senderID)
		if receiver == nil {
			log.Printf("No valid receiver found for send transaction from user %d", senderID)
			continue
		}
		receiverID := int(receiver.ID)

		if err := transactionService.CreateTransaction(senderID, -amount, "send"); err != nil {
			log.Printf("Failed to create send transaction for user %d: %v", senderID, err)
		} else {
			log.Printf("Created send transaction of %.2f from user %d to user %d", amount, senderID, receiverID)
		}

		if err := transactionService.CreateTransaction(receiverID, amount, "receive"); err != nil {
			log.Printf("Failed to create receive transaction for user %d: %v", receiverID, err)
		} else {
			log.Printf("Created receive transaction of %.2f for user %d from user %d", amount, receiverID, senderID)
		}
	}
}

func selectRandomUser(users []models.User, excludeID int) *models.User {
	var validUsers []models.User
	for _, user := range users {
		if int(user.ID) != excludeID {
			validUsers = append(validUsers, user)
		}
	}
	if len(validUsers) == 0 {
		return nil
	}
	return &validUsers[rand.Intn(len(validUsers))]
}
