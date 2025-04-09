package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"ptm/internal/repositories"
	"strconv"
	"strings"
	"sync"
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

func WarmUpBalanceCache() {
	userRepository := repositories.NewUserRepository()
	balanceRepository := repositories.NewBalanceRepository()
	page := uint(1)
	pageSize := uint(100)
	concurrency := 10
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for {
		users, err := userRepository.GetUsers(page, pageSize)
		if err != nil {
			log.Printf("Failed to fetch users on page %d: %v", page, err)
			break
		}
		if len(users) == 0 {
			break
		}

		for _, user := range users {
			wg.Add(1)
			sem <- struct{}{}
			go func(userID uint) {
				defer wg.Done()
				defer func() { <-sem }()
				balance, err := balanceRepository.GetBalance(userID)
				if err != nil {
					log.Printf("Failed to fetch balance for user %d: %v", userID, err)
					return
				}
				key := Key("balance", userID)
				if err := Set(key, balance, 10*time.Minute); err != nil {
					log.Printf("Failed to set balance for user %d: %v", userID, err)
				}
			}(user.ID)
		}

		page++
	}

	wg.Wait()
}

func Set(key string, value any, expiration ...time.Duration) error {
	var ttl time.Duration
	if len(expiration) > 0 {
		ttl = expiration[0]
	}
	err := redisClient.Set(context.Background(), key, fmt.Sprint(value), ttl).Err()
	if err != nil {
		log.Printf("Failed to set key %s in Redis: %v", key, err)
		return err
	}
	return nil
}

func Get(key string) (string, error) {
	val, err := redisClient.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		log.Printf("Key %s does not exist in Redis", key)
		return "", nil
	} else if err != nil {
		log.Printf("Failed to get key %s from Redis: %v", key, err)
		return "", err
	}
	return val, nil
}

func Key(parts ...any) string {
	s := make([]string, len(parts))
	for i, p := range parts {
		s[i] = fmt.Sprint(p)
	}
	return strings.Join(s, ":")
}

func Delete(key string) error {
	err := redisClient.Del(context.Background(), key).Err()
	if err != nil {
		log.Printf("Failed to delete key %s from Redis: %v", key, err)
		return err
	}
	return nil
}

func Exists(key string) (bool, error) {
	val, err := redisClient.Exists(context.Background(), key).Result()
	if err != nil {
		log.Printf("Failed to check if key %s exists in Redis: %v", key, err)
		return false, err
	}
	return val > 0, nil
}
