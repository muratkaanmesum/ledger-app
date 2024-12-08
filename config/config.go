package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	AppPort       string
	DatabaseURL   string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func InitConfig() Config {
	env := os.Getenv("APP_ENV") // Check if running locally or in Docker
	if env == "" {
		env = "local" // Default to local
	}

	var config Config
	config.AppPort = os.Getenv("APP_PORT")

	if env == "docker" {
		config.DatabaseURL = os.Getenv("DOCKER_DATABASE_URL")
		config.RedisAddr = os.Getenv("DOCKER_REDIS_ADDR")
	} else {
		config.DatabaseURL = os.Getenv("LOCAL_DATABASE_URL")
		config.RedisAddr = os.Getenv("LOCAL_REDIS_ADDR")
	}

	config.RedisPassword = os.Getenv("REDIS_PASSWORD")
	config.RedisDB = getEnvAsInt("REDIS_DB", 0)

	return config
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		var intValue int
		_, err := fmt.Sscanf(value, "%d", &intValue)
		if err != nil {
			log.Printf("Error parsing %s as int: %v\n", key, err)
			return defaultValue
		}
		return intValue
	}
	return defaultValue
}
