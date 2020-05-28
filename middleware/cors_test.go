package middleware

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"yee"
)

func TestCors(t *testing.T) {
	y := yee.New()
	y.Use(CorsWithConfig(CORSConfig{
		Origins:          []string{"*", "yearning.io"},
		AllowMethods:     []string{http.MethodGet},
		AllowHeaders:     []string{"Test"},
		AllowCredentials: false,
		ExposeHeaders:    nil,
		MaxAge:           0,
	}))

	y.GET("/", func(c yee.Context) error {
		return c.String(http.StatusOK, "test")
	})

	y.OPTIONS("/", func(c yee.Context) error {
		return c.String(http.StatusOK, "")
	})

	t.Run("http_get", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		y.ServeHTTP(rec, req)
		assert := assert.New(t)
		assert.Equal("test", rec.Body.String())
		assert.Equal("*", rec.Header().Get(yee.HeaderAccessControlAllowOrigin))
	})

	t.Run("http_option", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		rec := httptest.NewRecorder()
		y.ServeHTTP(rec, req)
		assert := assert.New(t)
		assert.Equal(http.MethodGet, rec.Header().Get(yee.HeaderAccessControlAllowMethods))
		assert.Equal("Test", rec.Header().Get(yee.HeaderAccessControlAllowHeaders))
	})
}
