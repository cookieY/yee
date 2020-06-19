package middleware

import (
	"crypto/subtle"
	"errors"
	"github.com/cookieY/yee"
	"net/http"
	"strings"
	"time"
)

type CSRFConfig struct {
	TokenLength    uint8
	TokenLookup    string
	Key            string
	CookieName     string
	CookieDomain   string
	CookiePath     string
	CookieMaxAge   int
	CookieSecure   bool
	CookieHttpOnly bool
}

type csrfTokenCreator func(yee.Context) (string, error)

var CSRFDefaultConfig = CSRFConfig{
	TokenLength:  16,
	TokenLookup:  "header:" + yee.HeaderXCSRFToken,
	Key:          "csrf",
	CookieName:   "_csrf",
	CookieMaxAge: 28800,
}

func Csrf() yee.HandlerFunc {
	return CSRFWithConfig(CSRFDefaultConfig)
}

func CSRFWithConfig(config CSRFConfig) yee.HandlerFunc {

	if config.TokenLength == 0 {
		config.TokenLength = CSRFDefaultConfig.TokenLength
	}

	if config.TokenLookup == "" {
		config.TokenLookup = CSRFDefaultConfig.TokenLookup
	}

	if config.Key == "" {
		config.Key = CSRFDefaultConfig.Key
	}

	if config.CookieName == "" {
		config.CookieName = CSRFDefaultConfig.CookieName
	}

	if config.CookieMaxAge == 0 {
		config.CookieMaxAge = CSRFDefaultConfig.CookieMaxAge
	}

	proc := strings.Split(config.TokenLookup, ":")

	creator := csrfTokenFromHeader(proc[1])

	switch proc[0] {
	case "query":
		creator = csrfTokenFromQuery(proc[1])
	case "form":
		creator = csrfTokenFromForm(proc[1])
	}

	return func(context yee.Context) (err error) {

		// we fetch cookie from this request
		// if cookie haven`t token info
		// we need generate the token and create a new cookie
		// otherwise reuse token

		k, err := context.Cookie(config.CookieName)
		token := ""
		if err != nil {
			token = yee.RandomString(config.TokenLength)
		} else {
			token = k.Value
		}

		switch context.Request().Method {
		case http.MethodGet, http.MethodTrace, http.MethodOptions, http.MethodHead:
		default:
			clientToken, e := creator(context)

			if e != nil {
				return context.ServerError(http.StatusBadRequest, e.Error())
			}
			if !validateCSRFToken(token, clientToken) {
				return context.ServerError(http.StatusForbidden, "invalid csrf token")
			}
		}

		nCookie := new(http.Cookie)
		nCookie.Name = config.CookieName
		nCookie.Value = token
		if config.CookiePath != "" {
			nCookie.Path = config.CookiePath
		}
		if config.CookieDomain != "" {
			nCookie.Domain = config.CookieDomain
		}
		nCookie.Expires = time.Now().Add(time.Duration(config.CookieMaxAge) * time.Second)
		nCookie.Secure = config.CookieSecure
		nCookie.HttpOnly = config.CookieHttpOnly
		context.SetCookie(nCookie)

		context.Put(config.Key, token)
		context.SetHeader(yee.HeaderVary, yee.HeaderCookie)
		context.Next()
		return
	}
}

func csrfTokenFromHeader(header string) csrfTokenCreator {
	return func(context yee.Context) (string, error) {
		token := context.GetHeader(header)
		if token == "" {
			return "", errors.New("missing csrf token in the header string")
		}
		return token, nil
	}
}

func csrfTokenFromQuery(param string) csrfTokenCreator {
	return func(context yee.Context) (string, error) {
		token := context.QueryParam(param)
		if token == "" {
			return "", errors.New("missing csrf token in the query string")
		}
		return token, nil
	}
}

func csrfTokenFromForm(param string) csrfTokenCreator {
	return func(context yee.Context) (string, error) {
		token := context.FormValue(param)
		if token == "" {
			return "", errors.New("missing csrf token in the form string")
		}
		return token, nil
	}
}
func validateCSRFToken(token, clientToken string) bool {
	return subtle.ConstantTimeCompare([]byte(token), []byte(clientToken)) == 1
}
