package middleware

import (
	"github.com/cookieY/Yee"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogger(t *testing.T) {
	y := yee.New()
	y.Use(Logger())
	y.GET("/", func(context yee.Context) error {
		context.Logger().Critical("哈哈哈哈")
		return context.String(http.StatusOK, "ok")
	})
	t.Run("http_get", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		y.ServeHTTP(rec, req)
		assert := assert.New(t)
		assert.Equal("ok", rec.Body.String())
		assert.Equal(http.StatusOK, rec.Code)
	})
}
