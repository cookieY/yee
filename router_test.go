package yee

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterParam(t *testing.T) {

	y := New()
	y.GET("/hello/k/:b", func(c Context) error {
		return c.String(http.StatusOK, c.Params("b"))
	})
	req := httptest.NewRequest(http.MethodGet, "/hello/k/yee", nil)
	rec := httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, "yee", rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRouterParams(t *testing.T) {

	y := New()
	y.GET("/hello/k/:b/p/:j", func(c Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("%s-%s", c.Params("b"), c.Params("j")))
	})
	req := httptest.NewRequest(http.MethodGet, "/hello/k/yee/p/henry", nil)
	rec := httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, "yee-henry", rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRouterStaticPath(t *testing.T) {

	y := New()
	y.GET("/hello/k/*assets", func(c Context) error {
		return c.String(http.StatusOK, c.Params("assets"))
	})
	req := httptest.NewRequest(http.MethodGet, "/hello/k/assets/1.js", nil)
	rec := httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, "/assets/1.js", rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRouterMultiRoute(t *testing.T) {

	y := New()
	y.GET("/hello/k/*k", func(c Context) error {
		return c.String(http.StatusOK, c.Params("k"))
	})
	y.GET("/hello/k", func(c Context) error {
		return c.String(http.StatusOK, "is_ok")
	})
	req := httptest.NewRequest(http.MethodGet, "/hello/k/yee/version/route", nil)
	rec := httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, "/yee/version/route", rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/hello/k", nil)
	rec = httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, "is_ok", rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)

}

func TestRouterQueryParam(t *testing.T) {

	y := New()
	y.GET("/hello/query", func(c Context) error {
		return c.String(http.StatusOK, c.QueryParam("query"))
	})
	req := httptest.NewRequest(http.MethodGet, "/hello/query?query=henry", nil)
	rec := httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, "henry", rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)

}

type testCase struct {
	uri    string
	expect string
}

func TestRouterMixin(t *testing.T) {
	y := New()
	y.GET("/pay", func(c Context) error {
		return c.String(http.StatusOK, "pay")
	})
	y.GET("/pay/add", func(c Context) error {
		return c.String(http.StatusOK, c.QueryParam("person"))
	})
	y.GET("/pay/add/:id", func(c Context) error {
		return c.String(http.StatusOK, c.Params("id"))
	})
	y.GET("/pay/add/:id/:store", func(c Context) error {
		return c.String(http.StatusOK, c.Params("id")+c.Params("store"))
	})
	y.GET("/pay/dew", func(c Context) error {
		return c.String(http.StatusOK, "dew")
	})
	y.GET("/pay/dew/*account", func(c Context) error {
		return c.String(http.StatusOK, c.Params("account"))
	})

	c := []testCase{
		{
			uri:    "/pay",
			expect: "pay",
		},
		{
			uri:    "/pay/add?person=henry",
			expect: "henry",
		},
		{
			uri:    "/pay/add/1",
			expect: "1",
		},
		{
			uri:    "/pay/add/1/a",
			expect: "1a",
		},
		{
			uri:    "/pay/dew",
			expect: "dew",
		},
		{
			uri:    "/pay/dew/account/css/1.css",
			expect: "/account/css/1.css",
		},
	}

	for _, i := range c {
		req := httptest.NewRequest(http.MethodGet, i.uri, nil)
		rec := httptest.NewRecorder()
		y.ServeHTTP(rec, req)
		assert.Equal(t, i.expect, rec.Body.String())
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

// If you want to test routing performance, You can use benchmark_test to get it

// --- testing any method

func testRestfulApi() RestfulApi {

	var api RestfulApi

	api.Get = func(c Context) (err error) {
		return c.String(http.StatusOK, "get")
	}

	api.Post = func(c Context) (err error) {
		return c.String(http.StatusOK, "post")
	}

	api.Delete = func(c Context) (err error) {
		return c.String(http.StatusOK, "delete")
	}

	api.Put = func(c Context) (err error) {
		return c.String(http.StatusOK, "put")
	}

	return api
}

func userUpdate(c Context) (err error) {
	return c.String(http.StatusOK, "updated")
}

func userFetch(c Context) (err error) {
	return c.String(http.StatusOK, "get it")
}

func test2RestfulApi() RestfulApi {
	return RestfulApi{
		Get:  userFetch,
		Post: userUpdate,
	}
}

func TestAnyMethod(t *testing.T) {

	y := New()

	y.Restful("/", testRestfulApi())
	y.Restful("/user", test2RestfulApi())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, "get", rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodPost, "/", nil)
	rec = httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, "post", rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodPut, "/", nil)
	rec = httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, "put", rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodDelete, "/", nil)
	rec = httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, "delete", rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/user", nil)
	rec = httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, "get it", rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodPost, "/user", nil)
	rec = httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, "updated", rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodPut, "/user", nil)
	rec = httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	req = httptest.NewRequest(http.MethodDelete, "/user", nil)
	rec = httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
