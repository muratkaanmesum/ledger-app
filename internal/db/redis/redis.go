package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
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

func AppendEventToStream(stream string, data map[string]interface{}) error {
	return redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: data,
	}).Err()
}

func ReadEventsFromStream(stream string) ([]redis.XMessage, error) {
	return redisClient.XRange(ctx, stream, "-", "+").Result()
}

func SetJSON[T any](key string, value T, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return redisClient.Set(context.Background(), key, data, expiration).Err()
}

func GetJSON[T any](key string) (T, error) {
	var result T
	val, err := redisClient.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		return result, nil
	}
	if err != nil {
		return result, err
	}
	err = json.Unmarshal([]byte(val), &result)
	return result, err
}

func getConversionRate(defaultCurrency, targetCurrency string) (float64, error) {
	type response struct {
		Result          string             `json:"result"`
		ConversionRates map[string]float64 `json:"conversion_rates"`
	}

	cacheKey := fmt.Sprintf("exchange_rate:%s:%s", defaultCurrency, targetCurrency)
	cachedRate, err := GetJSON[float64](cacheKey)
	if err == nil && cachedRate != 0 {
		return cachedRate, nil
	}

	client := resty.New()
	apiKey := os.Getenv("EXCHANGE_API_KEY")
	apiURL := "https://v6.exchangerate-api.com/v6/" + apiKey + "/latest/" + defaultCurrency

	resp, err := client.R().
		SetResult(&response{}).
		Get(apiURL)

	if err != nil {
		return 0, err
	}

	result := resp.Result().(*response)
	rate, ok := result.ConversionRates[targetCurrency]
	if !ok {
		return 0, fmt.Errorf("conversion rate for %s not found", targetCurrency)
	}

	_ = SetJSON(cacheKey, rate, time.Hour)
	return rate, nil
}
