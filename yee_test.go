package yee

import (
	"net/http"
	"testing"
)

func indexHandle(c Context) (err error) {
	return c.JSON(http.StatusOK, "ok")
}

func addRouter(y *Core) {
	y.GET("/", indexHandle)
}

func TestYee(t *testing.T) {
	y := New()
	addRouter(y)
	y.Run(":9999")
}

func TestRestApi(t *testing.T) {
	y := New()
	y.Restful("/", RestfulAPI{
		Get: func(c Context) (err error) {
			return c.String(http.StatusOK, "updated")
		},
		Post: func(c Context) (err error) {
			return c.String(http.StatusOK, "get it")
		},
	})
}

func TestDownload(t *testing.T)  {
	y := New()
	y.GET("/", func(c Context) (err error) {
		return c.File("args.go")
	})
	y.Run(":9999")
}