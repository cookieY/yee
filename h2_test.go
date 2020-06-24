package yee

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"testing"

	"golang.org/x/net/http2"
)

func TestH2(T *testing.T) {
	y := New()
	y.GET("/", func(c Context) error {
		return c.String(http.StatusOK, "ok")
	})
	y.RunH2C(":9999")
}

func TestH2cClient(t *testing.T) {
	client := http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}

	resp, err := client.Get("http://localhost:9999")
	if err != nil {
		log.Fatalf("faild request: %s", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("read response failed: %s", err)
	}
	fmt.Printf("proto:%s\ncode %d: %s\n", resp.Proto, resp.StatusCode, string(body))
}
