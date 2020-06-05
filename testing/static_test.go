package testing

import (
	"log"
	"net/http"
	"testing"
	"yee"
	"yee/middleware"
)

func TestStatic (t *testing.T) {
	y := yee.New()
	y.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level:9}))
	y.Static("/","./h5")
	y.Start(":9999")

}
func TestStatic2 (t *testing.T) {

	// 设置静态目录
	fsh := http.FileServer(http.Dir("./h5"))
	http.Handle("/", http.StripPrefix("/", fsh))

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
