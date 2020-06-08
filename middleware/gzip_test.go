package middleware

import (
	"github.com/cookieY/Yee"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGzip(t *testing.T) {
	y := yee.New()
	y.Use(Logger(), GzipWithConfig(GzipConfig{Level: 9}))
	y.Static("/", "../testing/dist/assets")
	t.Run("http_get", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/js/app.d2880701.js", nil)
		req.Header.Add(yee.HeaderAcceptEncoding, "gzip")
		rec := httptest.NewRecorder()
		y.ServeHTTP(rec, req)
		assert2 := assert.New(t)
		assert2.Equal(http.StatusOK, rec.Code)
	})
}
