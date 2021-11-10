package yee

import (
	"fmt"
	pb "github.com/cookieY/yee/test"
	"net/http"
	"testing"
)

const addr = "https://www.henry.com:9999/hello"

func TestH3Server(t *testing.T) {
	y := New()
	y.SetLogLevel(5)
	y.POST("/hello", func(c Context) (err error) {
		u := new(pb.Svr)
		if err := c.Bind(u); err != nil {
			c.Logger().Error(err.Error())
			return err
		}
		c.Logger().Debugf("svr get client data: %s",u.Project)
		svr := pb.Svr{Cloud: "hi"}
		return c.ProtoBuf(http.StatusOK, &svr)
	})
	y.RunH3(":9999", "henry.com+4.pem", "henry.com+4-key.pem")
}

func TestH3SvrIndex(t *testing.T)  {
	y := New()
	y.SetLogLevel(5)
	y.GET("/", func(c Context) (err error) {

		return c.JSON(http.StatusOK, "hello")
	})
	y.RunH3(":443", "henry.com+4.pem", "henry.com+4-key.pem")
}

func TestNewH3Client(t *testing.T) {
	cs := NewH3Client(&CConfig{
		Addr:               addr,
		InsecureSkipVerify: true,
	})
	rsp := new(pb.Svr)
	cs.Post(&pb.Svr{Project: "henry"}, rsp)
	fmt.Println(rsp.Cloud)
}

func TestNewProtoc3(t *testing.T)  {
	y := New()
	y.SetLogLevel(5)
	y.POST("/hello", func(c Context) (err error) {
		u := new(pb.Svr)
		if err := c.Bind(u); err != nil {
			c.Logger().Error(err.Error())
			return err
		}
		c.Logger().Debugf("svr get client data: %s",u.Project)
		svr := pb.Svr{Cloud: "hi"}
		return c.ProtoBuf(http.StatusOK, &svr)
	})
	y.Run(":9999")
}
