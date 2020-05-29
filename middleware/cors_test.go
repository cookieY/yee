package middleware

import (
	"fmt"
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

	y.GET("/ok", func(c yee.Context) error {
		return c.String(http.StatusOK, "test")
	})

	y.OPTIONS("/ok", func(context yee.Context) (err error) {
		return err
	})

	t.Run("http_get", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		rec := httptest.NewRecorder()
		y.ServeHTTP(rec, req)
		assert := assert.New(t)
		assert.Equal("test", rec.Body.String())
		assert.Equal("*", rec.Header().Get(yee.HeaderAccessControlAllowOrigin))
	})

	t.Run("http_option", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/ok", nil)
		rec := httptest.NewRecorder()
		y.ServeHTTP(rec, req)
		fmt.Println(rec.Body.String())
		assert := assert.New(t)
		assert.Equal(http.MethodGet, rec.Header().Get(yee.HeaderAccessControlAllowMethods))
		assert.Equal("Test", rec.Header().Get(yee.HeaderAccessControlAllowHeaders))
	})
}
