package middleware

import (
	"github.com/cookieY/yee"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID(t *testing.T) {
	y := yee.New()
	y.Use(RequestID())
	y.GET("/", func(context yee.Context) error {
		return context.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	y.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEqual(t, "", rec.Header().Get(yee.HeaderXRequestID))
}
