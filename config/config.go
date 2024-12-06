package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	AppPort     string
}

var AppConfig Config

func InitConfig() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to system environment variables.")
	}

	// Set config values
	AppConfig.DatabaseURL = getEnv("DATABASE_URL", "postgres://default:default@localhost:5432/defaultdb")
	AppConfig.AppPort = getEnv("APP_PORT", "8080")
}

// getEnv gets a key from the environment, or returns a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
