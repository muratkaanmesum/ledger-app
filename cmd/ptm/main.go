package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"os"
	"ptm/configs"
	"ptm/internal/db"
	"ptm/internal/db/redis"
	"ptm/internal/db/seeder"
	"ptm/internal/di"
	"ptm/internal/routes"
	"ptm/pkg/logger"
	"ptm/pkg/validator"
)

func main() {
	config.InitConfig()
	err := logger.InitLogger()
	if err != nil {
		return
	}
	err = godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables only")
	}

	db.InitDB()
	redis.InitRedis()
	di.InitDiContainer()

	e := echo.New()
	e.Validator = validator.New()

	fmt.Println(os.Getenv("ENVİRONMENT"))
	if os.Getenv("APP_ENV") == "development" {
		seeder.SeedUsers()
	}
	routes.InitRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}
