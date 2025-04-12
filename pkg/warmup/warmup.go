package warmup

import (
	"log"
	"ptm/internal/db/redis"
	"ptm/internal/repositories"
	"sync"
	"time"
)

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
				key := redis.Key("balance", userID)
				if err := redis.Set(key, balance, 10*time.Minute); err != nil {
					log.Printf("Failed to set balance for user %d: %v", userID, err)
				}
			}(user.ID)
		}

		page++
	}

	wg.Wait()
}
