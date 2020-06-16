package middleware

import (
	"github.com/cookieY/yee"
	"net/http"
	"sync"
	"time"
)

type RateLimitConfig struct {
	Time    time.Duration
	Rate    int
	lock    *sync.Mutex
	numbers int
}

var DefaultRateLimit = RateLimitConfig{
	Time: 1 * time.Second,
	Rate: 5,
}

func RateLimit() yee.HandlerFunc {
	return RateLimitWithConfig(DefaultRateLimit)
}

func RateLimitWithConfig(config RateLimitConfig) yee.HandlerFunc {

	if config.Time == 0 {
		config.Time = DefaultRateLimit.Time
	}

	if config.Rate == 0 {
		config.Rate = DefaultRateLimit.Rate
	}

	config.lock = new(sync.Mutex)

	go timer(&config)

	return yee.HandlerFunc{
		Func: func(context yee.Context) (err error) {
			if config.numbers >= config.Rate {
				context.ServerError(http.StatusTooManyRequests, "too many requests")
			}
			config.lock.Lock()
			config.numbers++
			defer config.lock.Unlock()
			return
		},
		IsMiddleware: true,
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
