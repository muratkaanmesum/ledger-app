package seeder

import (
	"log"
	"math/rand"
	"ptm/internal/db"
	"ptm/internal/models"
	"ptm/internal/services"
)

type userSeedData struct {
	Username string
	Role     string
	Email    string
	Password string
}

func SeedUsers() {
	users := []userSeedData{
		{Username: "John Doe", Role: "admin", Email: "test@gmail.com", Password: "123"},
		{Username: "Jane Smith", Role: "user", Email: "test1@gmail.com", Password: "asd"},
		{Username: "Alice Johnson", Role: "user", Email: "test2@gail.com", Password: "zxc"},
	}

	for _, user := range users {
		var existingUser models.User
		if err := db.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
			log.Printf("User with username %s already exists. Skipping seed.", user.Username)
			continue
		}
		createdUser, err := models.NewUser(user.Username, user.Email, user.Password, user.Role)
		if err != nil {
			log.Printf("Failed to create user %s: %v", user.Username, err)
			continue
		}
		if err := db.DB.Create(&createdUser).Error; err != nil {
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
		senderID := uint(sender.ID)

		depositAmount := float64(rand.Intn(500) + 50)
		if _, err := transactionService.CreateTransaction(senderID, senderID, depositAmount, "deposit"); err != nil {
			log.Printf("Failed to create deposit transaction for user %d: %v", senderID, err)
		} else {
			log.Printf("Created deposit transaction of %.2f for user %d", depositAmount, senderID)
		}

		withdrawAmount := float64(rand.Intn(200) + 10)
		if _, err := transactionService.CreateTransaction(senderID, senderID, -withdrawAmount, "withdraw"); err != nil {
			log.Printf("Failed to create withdraw transaction for user %d: %v", senderID, err)
		} else {
			log.Printf("Created withdraw transaction of %.2f for user %d", withdrawAmount, senderID)
		}

		amount := float64(rand.Intn(200) + 20)
		receiver := selectRandomUser(users, int(senderID))
		if receiver == nil {
			log.Printf("No valid receiver found for send transaction from user %d", senderID)
			continue
		}
		receiverID := uint(receiver.ID)

		if _, err := transactionService.CreateTransaction(senderID, receiverID, -amount, models.TransactionTypeDebit); err != nil {
			log.Printf("Failed to create send transaction for user %d: %v", senderID, err)
		} else {
			log.Printf("Created send transaction of %.2f from user %d to user %d", amount, senderID, receiverID)
		}

		if _, err := transactionService.CreateTransaction(receiverID, senderID, amount, models.TransactionTypeDebit); err != nil {
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
