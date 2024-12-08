package db

import (
	"log"
	"ptm/db/redis"
	"ptm/models"
	"strconv"
)

func SeedUsers() {
	users := []models.User{
		{Name: "John Doe", Role: "admin"},
		{Name: "Jane Smith", Role: "normal"},
		{Name: "Alice Johnson", Role: "normal"},
	}

	for _, user := range users {
		var existingUser models.User
		if err := DB.Where("name = ?", user.Name).First(&existingUser).Error; err == nil {
			log.Printf("User with name %s already exists. Skipping seed.", user.Name)
			continue
		}

		if err := DB.Create(&user).Error; err != nil {
			log.Printf("Failed to seed user %s: %v", user.Name, err)
			continue
		} else {
			log.Printf("Seeded user: %s with role %s", user.Name, user.Role)
		}

		if err := redis.Set(strconv.Itoa(int(user.ID)), "0", 0); err != nil {
			log.Printf("Failed to set Redis balance for user %s: %v", user.Name, err)
		} else {
			log.Printf("Initialized Redis balance for user %s", user.Name)
		}
	}
}
