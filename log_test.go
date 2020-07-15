package yee

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogger_LogWrite(t *testing.T) {

	y := New()
	y.SetLogLevel(3)

	y.POST("/hello/k/:b", func(c Context) error {
		c.Logger().Critical("critical")
		c.Logger().Error("error")
		c.Logger().Warn("warn")
		c.Logger().Info("info")
		c.Logger().Debug("debug")
		return c.String(http.StatusOK, c.Params("b"))
	})
	t.Run("http_get", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/hello/k/henry", nil)
		rec := httptest.NewRecorder()
		y.ServeHTTP(rec, req)
		tx := assert.New(t)
		tx.Equal("henry", rec.Body.String())
		//assert.Equal("*", rec.Header().Get(yee.HeaderAccessControlAllowOrigin))
	})
}
