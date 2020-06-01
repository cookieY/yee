package middleware

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"yee"
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
	r := y.Group("/", Cors(), JWTWithConfig(JwtConfig{SigningKey: []byte("dbcjqheupqjsuwsm")}),Logger())
	r.GET("/test", func(context yee.Context) error {
		return context.String(http.StatusOK, "is_ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	a := assert.New(t)
	y.ServeHTTP(rec, req)
	a.Equal("missing or malformed jwt", rec.Body.String())
	a.Equal(400, rec.Code)
	a.Equal("*", rec.Header().Get(yee.HeaderAccessControlAllowOrigin))
}
