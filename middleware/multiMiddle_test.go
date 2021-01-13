package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cookieY/yee"
	"github.com/stretchr/testify/assert"
)

func TestMultiMiddle(t *testing.T) {
	y := yee.New()
	y.Use(Cors())
	y.Use(JWTWithConfig(JwtConfig{SigningKey: []byte("dbcjqheupqjsuwsm")}))
	y.GET("/", func(context yee.Context) error {
		return context.String(http.StatusOK, "is_ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	a := assert.New(t)
	y.ServeHTTP(rec, req)
	a.Equal("\"missing or malformed jwt\"\n", rec.Body.String())
	a.Equal(400, rec.Code)
	a.Equal("*", rec.Header().Get(yee.HeaderAccessControlAllowOrigin))
}

func TestMultiGroup(t *testing.T) {
	y := yee.C()
	r := y.Group("/", Cors(), CustomerMiddleware())
	r.GET("/test", func(context yee.Context) error {
		return context.String(http.StatusOK, "is_ok")
	})
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	a := assert.New(t)
	y.ServeHTTP(rec, req)
	a.Equal("非法越权操作！", rec.Body.String())
	a.Equal("*", rec.Header().Get(yee.HeaderAccessControlAllowOrigin))
	a.Equal(403, rec.Code)
}

func CustomerMiddleware() yee.HandlerFunc {
	return func(c yee.Context) (err error) {
		if c.QueryParam("test") == "y" {
			c.Next()
			return
		}

		return c.ServerError(http.StatusForbidden, "非法越权操作！")
	}
}
