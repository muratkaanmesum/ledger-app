package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	models2 "ptm/internal/models"
)

var DB *gorm.DB

func InitDB() {
	var err error
	var dsn = os.Getenv("DATABASE_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err := DB.AutoMigrate(
		&models2.User{},
		&models2.Transaction{},
		&models2.Balance{},
		&models2.AuditLog{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connected and migrated successfully!")
}
