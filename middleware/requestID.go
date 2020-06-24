package middleware

import (
	"strings"

	"github.com/cookieY/yee"
	"github.com/google/uuid"
)

// RequestIDConfig defines config of requestID middleware
type RequestIDConfig struct {
	generator func() string
}

// DefaultRequestIDConfig is the default config of requestID middleware
var DefaultRequestIDConfig = RequestIDConfig{
	generator: defaultGenerator,
}

func defaultGenerator() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

// RequestID is the default implementation of requestID middleware
func RequestID() yee.HandlerFunc {
	return RequestIDWithConfig(DefaultRequestIDConfig)
}

// RequestIDWithConfig is the custom implementation of requestID middleware
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
