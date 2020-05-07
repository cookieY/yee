package knocker

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
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

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.fetchRoute("GET", "/hello/*filepath")

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.pattern != "/hello/*filepath" {
		t.Fatal("should match /hello/*filepath")
	}

	if ps["filepath"] != "*filepath" {
		t.Fatal("name should be equal to 'geektutu'")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["filepath"])

}
