package main

import (
	"github.com/joho/godotenv"
	"log"
	"ptm/config"
	"ptm/db"
	"ptm/db/redis"
	"ptm/routes"
	"ptm/utils"

	"github.com/labstack/echo/v4"
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

	db.SeedUsers()

	routes.InitRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}
