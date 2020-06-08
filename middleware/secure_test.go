package middleware

import (
	"github.com/cookieY/yee"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecure(t *testing.T) {
	y := yee.New()
	y.Use(Secure())
	y.GET("/", func(context yee.Context) error {
		return context.String(http.StatusOK, "ok")
	})
	t.Run("http_get", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		y.ServeHTTP(rec, req)
		assert := assert.New(t)
		assert.Equal("ok", rec.Body.String())
		assert.Equal(http.StatusOK, rec.Code)
		assert.Equal("SAMEORIGIN", rec.Header().Get(yee.HeaderXFrameOptions))
		assert.Equal("1; mode=block", rec.Header().Get(yee.HeaderXXSSProtection))
		assert.Equal("nosniff", rec.Header().Get(yee.HeaderXContentTypeOptions))
	})
}
