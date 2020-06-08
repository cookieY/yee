package middleware

import (
	"github.com/cookieY/yee"
	"net/http"
	"strconv"
	"strings"
)

type CORSConfig struct {
	Origins          []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           int
}

var DefaultCORSConfig = CORSConfig{
	Origins:      []string{"*"},
	AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch, http.MethodHead, http.MethodOptions, http.MethodConnect, http.MethodTrace},
}

func Cors() yee.HandlerFunc {
	return CorsWithConfig(DefaultCORSConfig)
}

func CorsWithConfig(config CORSConfig) yee.HandlerFunc {

	if len(config.Origins) == 0 {
		config.Origins = DefaultCORSConfig.Origins
	}

	if len(config.AllowMethods) == 0 {
		config.AllowMethods = DefaultCORSConfig.AllowMethods
	}

	allowMethods := strings.Join(config.AllowMethods, ",")

	allowHeaders := strings.Join(config.AllowHeaders, ",")

	exposeHeaders := strings.Join(config.ExposeHeaders, ",")

	maxAge := strconv.Itoa(config.MaxAge)

	return yee.HandlerFunc{

		Func: func(c yee.Context) (err error) {

			localOrigin := c.GetHeader(yee.HeaderOrigin)

			allowOrigin := ""

			m := c.Request().Method

			for _, o := range config.Origins {
				if o == "*" && config.AllowCredentials {
					allowOrigin = localOrigin
					break
				}
				if o == "*" || o == localOrigin {
					allowOrigin = o
					break
				}
			}

			if m != http.MethodOptions {
				c.AddHeader(yee.HeaderVary, yee.HeaderOrigin)
				c.SetHeader(yee.HeaderAccessControlAllowOrigin, allowOrigin)
				if config.AllowCredentials {
					c.SetHeader(yee.HeaderAccessControlAllowCredentials, "true")
				}
				if exposeHeaders != "" {
					c.SetHeader(yee.HeaderAccessControlExposeHeaders, exposeHeaders)
				}
				c.Next()
				return
			}

			c.AddHeader(yee.HeaderVary, yee.HeaderOrigin)
			c.AddHeader(yee.HeaderVary, yee.HeaderAccessControlRequestMethod)
			c.AddHeader(yee.HeaderVary, yee.HeaderAccessControlRequestHeaders)
			c.SetHeader(yee.HeaderAccessControlAllowOrigin, allowOrigin)
			c.SetHeader(yee.HeaderAccessControlAllowMethods, allowMethods)
			if config.AllowCredentials {
				c.SetHeader(yee.HeaderAccessControlAllowCredentials, "true")
			}
			if allowHeaders != "" {
				c.SetHeader(yee.HeaderAccessControlAllowHeaders, allowHeaders)
			} else {
				h := c.GetHeader(yee.HeaderAccessControlRequestHeaders)
				if h != "" {
					c.SetHeader(yee.HeaderAccessControlAllowHeaders, h)
				}
			}
			if config.MaxAge > 0 {
				c.SetHeader(yee.HeaderAccessControlMaxAge, maxAge)
			}
			return
		},
		IsMiddleware: true,
	}
}
