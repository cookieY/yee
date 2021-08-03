package middleware

import (
	"errors"
	"fmt"
	"github.com/cookieY/yee"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func GenJwtToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "henry"
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	t, err := token.SignedString([]byte("dbcjqheupqjsuwsm"))
	if err != nil {
		return "", errors.New("JWT Generate Failure")
	}
	return t, nil
}

func JwtParse(c yee.Context) (string, string) {
	user := c.Get("auth").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["name"].(string),claims["name"].(string)
}

func SuperManageGroup() yee.HandlerFunc {
	return func(c yee.Context) (err error) {
		user, _ := JwtParse(c)
		if user == "henry" {
			return
		}
		return c.JSON(http.StatusForbidden, "非法越权操作！")
	}
}

func TestJwt(t *testing.T) {

	cases := []struct {
		Name     string
		Expected int
		IsSign   bool
		Expire   time.Duration
	}{
		{"not_token", 400, false, 0},
		{"test_is_ok", 200, true, 0},
		//{"test_is_expire", 401, true,time.Second * 1},
	}
	for _, i := range cases {
		t.Run(i.Name, func(t *testing.T) {
			y := yee.New()

			y.Use(JWTWithConfig(JwtConfig{SigningKey: []byte("dbcjqheupqjsuwsm")}))
			y.GET("/", func(context yee.Context) error {
				return context.String(http.StatusOK, "is_ok")
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			if i.IsSign {
				token, _ := GenJwtToken()
				req.Header.Set(yee.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
			}
			time.Sleep(i.Expire)
			y.ServeHTTP(rec, req)
			assert2 := assert.New(t)
			assert2.Equal(i.Expected, rec.Code)
		})
	}
}

func TestJwtExpire(t *testing.T) {
	y := yee.New()
	y.Use(Cors(), Secure())
	y.GET("/", func(context yee.Context) error {
		return context.String(http.StatusOK, "is_ok")
	})
	r := y.Group("/api", JWTWithConfig(JwtConfig{SigningKey: []byte("dbcjqheupqjsuwsm")}))
	k := r.Group("/k", SuperManageGroup())
	k.GET("/o", func(context yee.Context) (err error) {
		return context.JSON(http.StatusOK, "pk")
	})
	y.Run(":9999")

}
