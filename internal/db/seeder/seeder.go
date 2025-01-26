package seeder

import (
	"log"
	"math/rand"
	"ptm/internal/db"
	"ptm/internal/di"
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
		dbUser := db.DB.Where("username = ?", user.Username).First(&existingUser)
		if dbUser.Error == nil {
			log.Printf("User with username %s already exists. Skipping seed.", user.Username)
		}

		userService := di.Resolve[services.UserService]()
		createdUser, err := userService.RegisterUser(&models.User{
			Username:     user.Username,
			Role:         user.Role,
			Email:        user.Email,
			PasswordHash: user.Password,
		})

		if err != nil {
			log.Printf("Failed to register user %s. Error: %v. Skipping seed.", user.Username, err)
		}

		balanceService := di.Resolve[services.BalanceService]()
		existingBalance, err := balanceService.GetUserBalance(createdUser.ID)
		if existingBalance != nil {
			log.Printf("Balance already exists for user %s. Skipping balance creation.", createdUser.Username)
			continue
		}

		balance, err := balanceService.CreateBalance(createdUser)
		if err != nil {
			log.Printf("Failed to create balance for user %s. Error: %v.", createdUser.Username, err)
			continue
		}

		if balance != nil {
			log.Printf("Balance successfully created for user %s.", createdUser.Username)
		}
	}
}

func createSpecificTransactions(users []models.User, transactionService services.TransactionService) {
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

		if _, err := transactionService.CreateTransaction(senderID, receiverID, -amount, models.Debit); err != nil {
			log.Printf("Failed to create send transaction for user %d: %v", senderID, err)
		} else {
			log.Printf("Created send transaction of %.2f from user %d to user %d", amount, senderID, receiverID)
		}

		if _, err := transactionService.CreateTransaction(receiverID, senderID, amount, models.Debit); err != nil {
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
