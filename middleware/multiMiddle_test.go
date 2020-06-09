package middleware

import (
	"fmt"
	"github.com/cookieY/yee"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMultiMiddle(t *testing.T) {
	y := yee.New()
	y.Use(JWTWithConfig(JwtConfig{SigningKey: []byte("dbcjqheupqjsuwsm")}))
	y.Use(Cors())
	y.GET("/", func(context yee.Context) error {
		return context.String(http.StatusOK, "is_ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	a := assert.New(t)
	y.ServeHTTP(rec, req)
	a.Equal("missing or malformed jwt", rec.Body.String())
	a.Equal(400, rec.Code)
	a.Equal("*", rec.Header().Get(yee.HeaderAccessControlAllowOrigin))
}

func TestMultiGroup(t *testing.T) {
	y := yee.New()
	y.Use(Cors())
	r := y.Group("/",CustomerMiddleware())
	r.POST("/login", func(context yee.Context) error {
		return context.String(http.StatusOK, "is_ok")
	})
	//y.Start(":8000")
	req := httptest.NewRequest(http.MethodGet, "/test?test=33", nil)
	rec := httptest.NewRecorder()

	a := assert.New(t)
	y.ServeHTTP(rec, req)
	a.Equal("非法越权操作！", rec.Body.String())
	a.Equal("*", rec.Header().Get(yee.HeaderAccessControlAllowOrigin))
	a.Equal(403, rec.Code)
}

func CustomerMiddleware() yee.HandlerFunc {
	return yee.HandlerFunc{
		Func: func(c yee.Context) (err error) {
			if c.QueryParam("test") == "y" {
				fmt.Println("231")
				c.Next()
				return
			}
			c.ServerError(http.StatusForbidden, []byte("非法越权操作！"), true)
			return
		},
		IsMiddleware: true,
	}
}
