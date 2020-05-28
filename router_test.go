package yee

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

<<<<<<< HEAD
func TestDyncRouter(t *testing.T) {
=======
func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/",nil)
	r.addRoute("GET", "/hello/*filepath", nil)
	r.addRoute("GET", "/:l/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}
func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parserParts("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parserParts("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parserParts("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}
>>>>>>> 6f3cf6c96722a3e91bd032293e8a15c1475ddbbb

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
