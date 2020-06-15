package middleware

import (
	"encoding/base64"
	"errors"
	"github.com/cookieY/yee"
	"net/http"
	"strings"
)

type BasicAuthConfig struct {
	Validator fnValidator
	Realm     string
}

type fnValidator func([]byte) (error, bool)

const (
	basic = "basic"
)

func BasicAuth(fn fnValidator) yee.HandlerFunc {
	config := BasicAuthConfig{Validator: fn}
	config.Realm = "."
	return BasicAuthWithConfig(config)
}

func BasicAuthWithConfig(config BasicAuthConfig) yee.HandlerFunc {

	if config.Validator == nil {
		panic("yee: basic-auth middleware requires a validator function")
	}

	return yee.HandlerFunc{
		Func: func(context yee.Context) (err error) {
			_, decode := parserVerifyData(context)
			if err, verify := config.Validator(decode); err == nil && verify {
				return err
			}

			context.Response().Header().Set(yee.HeaderWWWAuthenticate, basic+" realm="+config.Realm)

			context.ServerError(http.StatusUnauthorized, "invalid basic auth token")

			return
		},
		IsMiddleware: true,
	}
}

func parserVerifyData(context yee.Context) (error, []byte) {
	var decode []byte
	res := context.Request()
	if res.Header.Get(yee.HeaderAuthorization) != "" {
		auth := strings.Split(res.Header.Get(yee.HeaderAuthorization), " ")
		if auth[0] == basic {
			decode, err := base64.StdEncoding.DecodeString(auth[1])
			if err != nil {
				return err, decode
			}
			return nil, decode
		}
		return errors.New("cannot get basic keyword"), decode
	}
	return errors.New("authorization header is empty"), decode
}
