package yee

import (
	"net/http"
	"testing"
)

func indexHandle(c Context) (err error) {
	return c.JSON(http.StatusOK,"ok")
}

func addRouter(y *Core)  {
	y.GET("/",indexHandle)
}

func TestYee(t *testing.T)  {
	y:= New()
	addRouter(y)
	y.Run(":9999")
}
