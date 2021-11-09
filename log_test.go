package yee

import (
	"github.com/cookieY/yee/logger"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestLogger_LogWrite(t *testing.T) {

	y := New()
	y.SetLogLevel(7)
	fs, _ := os.Open("1.log")
	y.SetLogOut(fs)
	y.POST("/hello/k/:b", func(c Context) error {
		c.Logger().Critical("critical")
		c.Logger().Error("error")
		c.Logger().Warn("warn")
		c.Logger().Info("info")
		c.Logger().Debug("debug")
		c.Logger().Criticalf("test:%v", 123)
		c.Logger().Errorf("test:%v", 123)
		c.Logger().Warnf("test:%v", 123)
		c.Logger().Infof("test:%v", 123)
		c.Logger().Debugf("test:%v", 123)
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

func BenchmarkLogger_LogWrite(b *testing.B) {
	l := logger.LogCreator()
	b.ReportAllocs()
	b.SetBytes(1024 * 1024)
	for i := 0; i < b.N; i++ {
		l.Critical("critical")
		l.Error("error")
		l.Warn("warn")
		l.Info("info")
		l.Debug("debug")
		l.Criticalf("test:%v", 123)
		l.Errorf("test:%v", 123)
		l.Warnf("test:%v", 123)
		l.Infof("test:%v", 123)
		l.Debugf("test:%v", 123)
	}
}
