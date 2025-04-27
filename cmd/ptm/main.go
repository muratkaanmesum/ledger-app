package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"os"
	"os/signal"
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
	"syscall"
	"time"
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
	
	go func() {
		if err := e.Start(":8080"); err != nil {
			log.Printf("Server shutting down: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down server...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Closing database connections...")
	db.CloseDB()
	
	log.Println("Closing Redis connections...")
	redis.CloseRedis()
	
	log.Println("Stopping scheduler...")
	scheduler.StopScheduler()
	
	log.Println("Server gracefully stopped")
}
