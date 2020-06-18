package yee

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testData = `{"id":1,"name":"Jon Snow"}`

type res struct {
	Id   int    `json:"id"`
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
	assert.Equal(t,testData+"\n",rec.Body.String())
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
		req := httptest.NewRequest(http.MethodPost, "/",nil)
		rec := httptest.NewRecorder()
		y.ServeHTTP(rec, req)
	}
}

