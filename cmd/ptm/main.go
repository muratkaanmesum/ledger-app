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
	"ptm/internal/monitoring"
	"ptm/internal/routes"
	"ptm/internal/scheduler"
	"ptm/pkg/counter"
	"ptm/pkg/logger"
	"ptm/pkg/validator"
	"ptm/pkg/warmup"
)

func main() {
	config.InitConfig()
	err := logger.InitLogger()
	if err != nil {
		return
	}
	err = godotenv.Load()
	if err != nil {
		fmt.Println("err is ", err)
		log.Println("No .env file found, using environment variables only")
	}

	db.InitDB()
	redis.InitRedis()
	warmup.WarmUpBalanceCache()
	monitoring.InitPrometheus()

	scheduler.InitScheduler()
	if err != nil {
		log.Fatalf("Failed to initialize tracing: %v", err)
	}
	di.InitDiContainer()
	counter.InitStats()

	e := echo.New()
	e.Validator = validator.New()

	if os.Getenv("APP_ENV") == "development" {
		seeder.SeedUsers()
	}
	routes.InitRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}
