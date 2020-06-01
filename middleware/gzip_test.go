package middleware

import (
	"net/http"
	"testing"
	"yee"
)

func TestGzip(t *testing.T) {
	y := yee.New()
	y.Use(Logger(),GzipWithConfig(GzipConfig{Level: 5}))
	y.GET("/", func(context yee.Context) error {
		return context.String(http.StatusOK, "okduiashiudhasiudhias")
	})
	y.Start(":9999")
	//t.Run("http_get", func(t *testing.T) {
	//	req := httptest.NewRequest(http.MethodGet, "/", nil)
	//	req.Header.Add(yee.HeaderAcceptEncoding,"gzip")
	//	rec := httptest.NewRecorder()
	//	y.ServeHTTP(rec, req)
	//	assert := assert.New(t)
	//	assert.Equal(http.StatusOK, rec.Code)
	//})
}
