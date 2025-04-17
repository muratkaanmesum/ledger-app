package exchange

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
	"ptm/internal/db/redis"
	"time"
)

func GetConversionRate(defaultCurrency, targetCurrency string) (float64, error) {
	type response struct {
		Result          string             `json:"result"`
		ConversionRates map[string]float64 `json:"conversion_rates"`
	}
	if defaultCurrency == targetCurrency {
		return 1, nil
	}

	cacheKey := fmt.Sprintf("exchange_rate:%s", defaultCurrency)
	cachedResp, err := redis.GetJSON[response](cacheKey)
	if err == nil {
		if rate, ok := cachedResp.ConversionRates[targetCurrency]; ok {
			return rate, nil
		}
		return 0, fmt.Errorf("conversion rate for %s not found in cached data", targetCurrency)
	}

	client := resty.New()
	apiKey := os.Getenv("EXCHANGE_API_KEY")
	apiURL := "https://v6.exchangerate-api.com/v6/" + apiKey + "/latest/" + defaultCurrency
	var resp = response{}
	_, err = client.R().
		SetResult(&resp).
		Get(apiURL)

	if err != nil {
		return 0, err
	}

	if resp.ConversionRates == nil {
		return 0, fmt.Errorf("no conversion rates available")
	}

	_ = redis.SetJSON(cacheKey, resp, time.Hour)

	rate, ok := resp.ConversionRates[targetCurrency]
	if !ok {
		return 0, fmt.Errorf("conversion rate for %s not found", targetCurrency)
	}

	return rate, nil
}
