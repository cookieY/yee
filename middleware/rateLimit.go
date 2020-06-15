package middleware

import "github.com/cookieY/yee"

type RateLimitConfig struct {
	Rate    int
	Seconds int
}

var DefaultRateLimitConfig = RateLimitConfig{
	Rate:    5,
	Seconds: 1,
}

func RateLimitWithConfig(r RateLimitConfig) yee.HandlerFunc {

	if r.Rate == 0 {
		r.Rate = DefaultRateLimitConfig.Rate
	}

	if r.Seconds == 0 {
		r.Rate = DefaultRateLimitConfig.Seconds
	}

	return yee.HandlerFunc{
		Func: func(context yee.Context) (err error) {
			return
		},
		IsMiddleware: true,
	}
}
