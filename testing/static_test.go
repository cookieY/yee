package testing

import (
	"net/http"
	"testing"
	"yee"
	"yee/middleware"
)

func RouterGroups (y *yee.Core) {
	y.GET("/", func(c yee.Context) error {
		return c.HTMLTml(http.StatusOK,"./dist/index.html")
	})
}

func TestStatic(t *testing.T) {
	y := yee.New()
	y.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{IsLogger:true}), middleware.GzipWithConfig(middleware.GzipConfig{Level: 9}))
	RouterGroups(y)
	//y.GET("/api/v1", func(c yee.Context) error {
	//	return c.String(200, "ok")
	//})
	y.Static("/assets", "./dist/assets")
	y.Start(":9999")

}