package middleware

import (
	"github.com/cookieY/yee"
)

type RequestIDConfig struct {
	generator func() string
}

var DefaultRequestIDConfig = RequestIDConfig{
	generator: defaultGenerator,
}

func defaultGenerator() string {
	return yee.RandomString(16)
}

func RequestID() yee.HandlerFunc {
	return RequestIDWithConfig(DefaultRequestIDConfig)
}

func RequestIDWithConfig(config RequestIDConfig) yee.HandlerFunc {

	if config.generator == nil {
		config.generator = DefaultRequestIDConfig.generator
	}
	return func(context yee.Context) (err error) {
		req := context.Request()
		res := context.Response()
		if req.Header.Get(yee.HeaderXRequestID) == "" {
			res.Header().Set(yee.HeaderXRequestID, config.generator())
		}
		return
	}
}
