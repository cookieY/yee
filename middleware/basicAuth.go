package middleware

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/cookieY/yee"
)

// BasicAuthConfig defines the config of basicAuth middleware
type BasicAuthConfig struct {
	Validator fnValidator
	Realm     string
}

type fnValidator func([]byte) (bool, error)

const (
	basic = "basic"
)

// BasicAuth is the default implementation BasicAuth middleware
func BasicAuth(fn fnValidator) yee.HandlerFunc {
	config := BasicAuthConfig{Validator: fn}
	config.Realm = "."
	return BasicAuthWithConfig(config)
}

// BasicAuthWithConfig is the custom implementation BasicAuth middleware
func BasicAuthWithConfig(config BasicAuthConfig) yee.HandlerFunc {

	if config.Validator == nil {
		panic("yee: basic-auth middleware requires a validator function")
	}

	return func(context yee.Context) (err error) {
		decode, _ := parserVerifyData(context)
		if verify, err := config.Validator(decode); err == nil && verify {
			return err
		}

		context.Response().Header().Set(yee.HeaderWWWAuthenticate, basic+" realm="+config.Realm)

		return context.ServerError(http.StatusUnauthorized, "invalid basic auth token")
	}
}

func parserVerifyData(context yee.Context) ([]byte, error) {
	var decode []byte
	res := context.Request()
	if res.Header.Get(yee.HeaderAuthorization) != "" {
		auth := strings.Split(res.Header.Get(yee.HeaderAuthorization), " ")
		if auth[0] == basic {
			decode, err := base64.StdEncoding.DecodeString(auth[1])
			if err != nil {
				return decode, err
			}
			return decode, nil
		}
		return decode, errors.New("cannot get basic keyword")
	}
	return decode, errors.New("authorization header is empty")
}
