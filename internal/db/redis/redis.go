package redis

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("Invalid REDIS_DB value: %v", err)
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully!")
}

func Set(key string, value string, expiration time.Duration) error {
	err := redisClient.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		log.Printf("Failed to set key %s in Redis: %v", key, err)
		return err
	}
	return nil
}

// Get retrieves a value by key from Redis
func Get(key string) (string, error) {
	val, err := redisClient.Get(context.Background(), key).Result()
	if err == redis.Nil {
		log.Printf("Key %s does not exist in Redis", key)
		return "", nil // Key does not exist
	} else if err != nil {
		log.Printf("Failed to get key %s from Redis: %v", key, err)
		return "", err
	}
	return val, nil
}

// Delete removes a key from Redis
func Delete(key string) error {
	err := redisClient.Del(context.Background(), key).Err()
	if err != nil {
		log.Printf("Failed to delete key %s from Redis: %v", key, err)
		return err
	}
	return nil
}

// Exists checks if a key exists in Redis
func Exists(key string) (bool, error) {
	val, err := redisClient.Exists(context.Background(), key).Result()
	if err != nil {
		log.Printf("Failed to check if key %s exists in Redis: %v", key, err)
		return false, err
	}
	return val > 0, nil
}
