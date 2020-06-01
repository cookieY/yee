package middleware

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"yee"
)

func TestRecovery(t *testing.T) {
	y := yee.New()
	y.Use(Recovery())
	y.GET("/", func(context yee.Context) error {
		var t error
		return context.String(http.StatusOK, t.Error())
	})
	t.Run("http_get", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		y.ServeHTTP(rec, req)
		assert := assert.New(t)
		assert.Equal("Internal Server Error", rec.Body.String())
		assert.Equal(http.StatusInternalServerError, rec.Code)
	})
}

