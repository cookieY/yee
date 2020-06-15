package middleware

import (
	"errors"
	"fmt"
	"github.com/cookieY/yee"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"reflect"
	"strings"
)

type JwtConfig struct {
	GetKey        string
	AuthScheme    string
	SigningKey    interface{}
	SigningMethod string
	TokenLookup   string
	Claims        jwt.Claims
	keyFunc       jwt.Keyfunc
	ErrorHandler  JWTErrorHandler
}

type jwtExtractor func(yee.Context) (string, error)

type JWTErrorHandler func(error) error

const (
	AlgorithmHS256 = "HS256"
)

var DefaultJwtConfig = JwtConfig{
	GetKey:        "auth",
	SigningMethod: AlgorithmHS256,
	AuthScheme:    "Bearer",
	TokenLookup:   "header:" + yee.HeaderAuthorization,
	Claims:        jwt.MapClaims{},
}

func JWTWithConfig(config JwtConfig) yee.HandlerFunc {
	if config.SigningKey == nil {
		panic("yee: jwt middleware requires signing key")
	}
	if config.SigningMethod == "" {
		config.SigningMethod = DefaultJwtConfig.SigningMethod
	}
	if config.GetKey == "" {
		config.GetKey = DefaultJwtConfig.GetKey
	}
	if config.AuthScheme == "" {
		config.AuthScheme = DefaultJwtConfig.AuthScheme
	}

	if config.Claims == nil {
		config.Claims = DefaultJwtConfig.Claims
	}

	if config.TokenLookup == "" {
		config.TokenLookup = DefaultJwtConfig.TokenLookup
	}

	config.keyFunc = func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != config.SigningMethod {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", token.Header["alg"])
		}
		return config.SigningKey, nil
	}

	parts := strings.Split(config.TokenLookup, ":")
	extractor := jwtFromHeader(parts[1], config.AuthScheme)
	return yee.HandlerFunc{
		Func: func(c yee.Context) (err error) {
			auth, err := extractor(c)
			if err != nil {
				c.ServerError(http.StatusBadRequest, err.Error())
				return err
			}
			token := new(jwt.Token)
			if _, ok := config.Claims.(jwt.MapClaims); ok {
				token, err = jwt.Parse(auth, config.keyFunc)
				if err != nil {
					c.ServerError(http.StatusUnauthorized, err.Error())
					return err
				}
			} else {
				t := reflect.ValueOf(config.Claims).Type().Elem()
				claims := reflect.New(t).Interface().(jwt.Claims)
				token, err = jwt.ParseWithClaims(auth, claims, config.keyFunc)
			}
			if err == nil && token.Valid {
				c.Put(config.GetKey, token)
				return
			}
			// bug fix
			// if  invalid or expired jwt,
			// we must intercept all handlers and return serverError
			c.ServerError(http.StatusUnauthorized, "invalid or expired jwt")
			return
		},
		IsMiddleware: true,
	}
}

func jwtFromHeader(header string, authScheme string) jwtExtractor {
	return func(c yee.Context) (string, error) {
		auth := c.Request().Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", errors.New("missing or malformed jwt")
	}
}
