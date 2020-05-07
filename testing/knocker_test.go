package testing

import (
	"fmt"
	"knocker"
	"log"
	"net/http"
	"testing"
	"time"
)

func onlyForV2() knocker.HandlerFunc {
	return func(c knocker.Context) error {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Status(500)
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", 500, c.Scheme(), time.Since(t))
		return nil
	}
}

func TestContext(t *testing.T) {
	r := knocker.New()
	//r.GET("/", func(c knocker.Context) error {
	//	return c.String(http.StatusOK, "<h1>Hello Gee</h1>")
	//})

	r.GET("/:test/xxx", func(c knocker.Context) (err error) {
		j := c.Params("test")
		x := c.Params("xx")
		fmt.Println(1)
		return c.JSON(http.StatusOK, fmt.Sprintf("%s-%s", j, x))
		//return c.JSON(http.StatusOK, map[string]interface{}{"name":"henry","age": 27})
	})

	r.GET("/b/:name", func(c knocker.Context) (err error) {
		j := c.Params("name")
		fmt.Println(2)
		return c.JSON(http.StatusOK, j)
		//return c.JSON(http.StatusOK, map[string]interface{}{"name":"henry","age": 27})
	})

	//r.Static("/", "dist")
	//v1 := r.Group("/v1")
	//v1.Use(middleware.Xss())
	//v1.GET("/fail", func(c knocker.Context) (err error) {
	//	return
	//	//return c.JSON(http.StatusOK, map[string]interface{}{"name":"henry","age": 27})
	//})
	_ = r.Start(":9999")
}
