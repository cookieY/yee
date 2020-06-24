package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/cookieY/yee"
)

// RateLimitConfig defines config of rateLimit middleware
type RateLimitConfig struct {
	Time    time.Duration
	Rate    int
	lock    *sync.Mutex
	numbers int
}

// DefaultRateLimit is the default config of rateLimit middleware
var DefaultRateLimit = RateLimitConfig{
	Time: 1 * time.Second,
	Rate: 5,
}

// RateLimit is the default implementation of rateLimit middleware
func RateLimit() yee.HandlerFunc {
	return RateLimitWithConfig(DefaultRateLimit)
}

// RateLimitWithConfig is the custom implementation of rateLimit middleware
func RateLimitWithConfig(config RateLimitConfig) yee.HandlerFunc {

	if config.Time == 0 {
		config.Time = DefaultRateLimit.Time
	}

	if config.Rate == 0 {
		config.Rate = DefaultRateLimit.Rate
	}

	config.lock = new(sync.Mutex)

	go timer(&config)

	return func(context yee.Context) (err error) {
		if config.numbers >= config.Rate {
			return context.ServerError(http.StatusTooManyRequests, "too many requests")
		}
		config.lock.Lock()
		config.numbers++
		defer config.lock.Unlock()
		return
	}
}

func timer(c *RateLimitConfig) {
	ticker := time.NewTicker(c.Time)
	for {
		<-ticker.C
		c.lock.Lock()
		c.numbers = 0
		c.lock.Unlock()
	}
}
