package yee

import (
	"crypto/tls"
	"fmt"
	pb "github.com/cookieY/yee/test"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"io/ioutil"
	"net/http"
	"runtime"
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
		c.Logger().Debugf("svr get client data: %s", u.Project)
		svr := pb.Svr{Cloud: "hi"}
		return c.ProtoBuf(http.StatusOK, &svr)
	})
	y.RunH3(":9999", "henry.com+4.pem", "henry.com+4-key.pem")
}

func TestH3SvrIndex(t *testing.T) {
	y := New()
	y.SetLogLevel(5)
	y.GET("/", func(c Context) (err error) {

		return c.JSON(http.StatusOK, "hello")
	})
	y.Run(":445")
}

func TestH2SvrIndex(t *testing.T) {
	y := New()
	y.SetLogLevel(5)
	y.GET("/", func(c Context) (err error) {

		return c.JSON(http.StatusOK, "hello")
	})
	y.RunTLS(":444", "henry.com+4.pem", "henry.com+4-key.pem")
}

var cs = http.Client{
	Transport: &http3.RoundTripper{
		TLSClientConfig: &tls.Config{},
		QuicConfig:      &quic.Config{},
	},
}

func BenchmarkH3SvrIndex(b *testing.B) {
	b.SetBytes(1024 * 1024)
	for i := 0; i < b.N; i++ {
		http.Get("http://127.0.0.1:445/")
	}
}

func BenchmarkH2SvrIndex(b *testing.B) {
	b.SetBytes(1024 * 1024)
	for i := 0; i < b.N; i++ {
		http.Get("https://127.0.0.1:444/")
	}
}

func TestRespProto(t *testing.T) {
	//cs := http.Client{
	//	Transport: &http3.RoundTripper{
	//		TLSClientConfig: &tls.Config{
	//			InsecureSkipVerify: true,
	//		},
	//		QuicConfig: &quic.Config{},
	//	},
	//}

	b, err := http.Get("https://127.0.0.1:444/")
	if err != nil {
		t.Error(err)
	}
	c, _ := ioutil.ReadAll(b.Body)
	fmt.Println(string(c))
	fmt.Println(b.Proto)
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

func TestNewProtoc3(t *testing.T) {
	y := New()
	y.SetLogLevel(5)
	y.POST("/hello", func(c Context) (err error) {
		u := new(pb.Svr)
		if err := c.Bind(u); err != nil {
			c.Logger().Error(err.Error())
			return err
		}
		c.Logger().Debugf("svr get client data: %s", u.Project)
		svr := pb.Svr{Cloud: "hi"}
		return c.ProtoBuf(http.StatusOK, &svr)
	})
	y.Run(":9999")
}

func BenchmarkH2vsH3(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.Run("http3", BenchmarkH3SvrIndex)
	b.Run("http2", BenchmarkH2SvrIndex)
}
