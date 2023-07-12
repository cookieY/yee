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

func TestDownload(t *testing.T) {
	y := New()
	y.GET("/", func(c Context) (err error) {
		return c.File("args.go")
	})
	y.Run(":9999")
}

func TestStatic(t *testing.T) {
	y := New()
	y.Static("/", "dist")
	//y.GET("/", func(c Context) error {
	//	return c.HTMLTpl(http.StatusOK, "./dist/index.html")
	//})
	y.Run(":9999")
}

////go:embed dist/*
//var f embed.FS
//
////go:embed dist/index.html
//var index string

//func TestPack(t *testing.T) {
//	y := New()
//	y.Pack("/front", f, "dist")
//	y.GET("/", func(c Context) error {
//		return c.HTML(http.StatusOK, index)
//	})
//	y.Run(":9999")
//}

const ver = `alt-svc: h3=":443"; ma=2592000,h3-29=":443"; ma=2592000,h3-Q050=":443"; ma=2592000,h3-Q046=":443"; ma=2592000,h3-Q043=":443"; ma=2592000,quic=":443"; ma=2592000; v="46,43"`

func TestH3(t *testing.T) {
	y := New()
	y.GET("/", func(c Context) (err error) {
		return c.JSON(http.StatusOK, "hello")
	})
	y.RunH3(":443", "henry.com+4.pem", "henry.com+4-key.pem")
}
