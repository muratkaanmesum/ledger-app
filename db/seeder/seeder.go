package seeder

import (
	"log"
	"math/rand"
	"ptm/db"
	"ptm/models"
	"ptm/services"
)

func SeedUsers() {
	users := []models.User{
		{Username: "John Doe", Role: "admin"},
		{Username: "Jane Smith", Role: "normal"},
		{Username: "Alice Johnson", Role: "normal"},
	}

	for _, user := range users {
		var existingUser models.User
		if err := db.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
			log.Printf("User with username %s already exists. Skipping seed.", user.Username)
			continue
		}

		if err := db.DB.Create(&user).Error; err != nil {
			log.Printf("Failed to seed user %s: %v", user.Username, err)
			continue
		}
	}
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
