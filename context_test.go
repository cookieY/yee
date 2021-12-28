package yee

import (
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testData = `{"id":1,"name":"Jon Snow"}`

type res struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestContextJSON(t *testing.T) {
	y := New()
	y.POST("/", func(c Context) (err error) {
		t := new(res)
		if err = c.Bind(&t); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, t)
	})
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(testData))
	req.Header.Set("Content-Type", MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, testData+"\n", rec.Body.String())
}

func TestContextForward(t *testing.T) {
	y := New()
	y.POST("/", func(c Context) (err error) {
		return c.JSON(http.StatusOK, c.RemoteIP())
	})
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(HeaderXForwardedFor, "<img> ")
	rec := httptest.NewRecorder()
	y.ServeHTTP(rec, req)
}

func TestContextString(t *testing.T) {
	y := New()
	y.POST("/", func(c Context) (err error) {
		return c.String(http.StatusOK, "hello")
	})
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, "hello", rec.Body.String())
}

func crashMiddleware() HandlerFunc {
	return func(c Context) (err error) {
		c.CrashWithStatus(http.StatusUnauthorized)
		return
	}
}
func sayMiddleware() HandlerFunc {
	return func(c Context) (err error) {
		log.Println("say")
		return
	}
}

func TestCrash(t *testing.T) {
	y := New()
	y.Use(crashMiddleware())
	y.Use(sayMiddleware())
	y.GET("/", func(c Context) (err error) {
		return c.String(http.StatusOK, "hello")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Body)
}

func TestRedirect(t *testing.T) {
	y := New()
	y.GET("/", func(c Context) (err error) {
		return c.Redirect(http.StatusMovedPermanently, "/get")
	})
	y.GET("/get", func(c Context) (err error) {
		return c.String(http.StatusOK, "hello")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	y.ServeHTTP(rec, req)
}

func BenchmarkAllocJSON(b *testing.B) {
	y := New()
	y.POST("/", func(c Context) (err error) {
		tl := new(res)
		if err = c.Bind(&tl); err != nil {

			return err
		}
		return c.JSON(http.StatusOK, tl)
	})
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(testData))
		req.Header.Set("Content-Type", MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		y.ServeHTTP(rec, req)
	}
}

func BenchmarkAllocString(b *testing.B) {
	y := New()
	y.POST("/", func(c Context) (err error) {
		return c.String(http.StatusOK, "ok")
	})
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		y.ServeHTTP(rec, req)
	}
}
