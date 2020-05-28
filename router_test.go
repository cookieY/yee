package yee

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDyncRouter(t *testing.T) {

	y := New()



	y.GET("/hello/k/:name", func(c Context) error {
		return c.String(http.StatusOK, c.Params("name"))
	})

	y.GET("/hello/k/b", func(c Context) error {
		return c.String(http.StatusOK, "ok")
	})

	t.Run("http_get", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/hello/k/b2312312321", nil)
		rec := httptest.NewRecorder()
		y.ServeHTTP(rec, req)
		fmt.Println(rec.Body.String())
		//assert.Equal("test", rec.Body.String())
		//assert.Equal("*", rec.Header().Get(yee.HeaderAccessControlAllowOrigin))
	})
}
