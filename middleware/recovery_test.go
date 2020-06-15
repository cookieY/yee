package middleware

import (
	"github.com/cookieY/yee"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecovery(t *testing.T) {
	y := yee.InitCore()
	y.Use(Recovery())
	y.GET("/y", func(context yee.Context) error {
		names := []string{"geektutu"}
		return context.String(http.StatusOK, names[100])
	})

	t.Run("http_get", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/y", nil)
		rec := httptest.NewRecorder()
		y.ServeHTTP(rec, req)
		assert := assert.New(t)
		assert.Equal("Internal Server Error", rec.Body.String())
		assert.Equal(http.StatusInternalServerError, rec.Code)
	})
}
