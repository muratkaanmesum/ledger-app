package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"ptm/configs"
	"ptm/internal/db"
	"ptm/internal/db/redis"
	"ptm/internal/db/seeder"
	"ptm/internal/routes"
	"ptm/internal/utils"
	"ptm/internal/utils/validator"
)

func main() {
	config.InitConfig()
	utils.InitLogger()
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables only")
	}

	db.InitDB()
	redis.InitRedis()

	e := echo.New()
	e.Validator = validator.New()
	seeder.SeedUsers()

	routes.InitRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}
