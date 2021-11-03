package yee

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"io/ioutil"
	"net/http"
)

type transport struct {
	addr               string
	insecureSkipVerify bool
	logger             *logger
	tripper            *http3.RoundTripper
	c                  *http.Client
}

type CConfig struct {
	Addr               string
	InsecureSkipVerify bool
}

func NewH3Client(c *CConfig) *transport {
	tripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		QuicConfig: &quic.Config{},
	}
	return &transport{
		addr:               c.Addr,
		insecureSkipVerify: c.InsecureSkipVerify,
		logger:             LogCreator(),
		c: &http.Client{
			Transport: tripper,
		},
		tripper: tripper,
	}
}

func (t *transport) Post(payload proto.Message, recv proto.Message) {
	p, err := proto.Marshal(payload)
	if err != nil {
		t.logger.Critical(err.Error())
		return
	}
	fmt.Println(t.addr)
	rsp, err := t.c.Post(t.addr, MIMEApplicationProtobuf, bytes.NewReader(p))
	if err != nil {
		t.logger.Critical(err.Error())
		return
	}
	b, err := ioutil.ReadAll(rsp.Body)
	if rsp.StatusCode == 200 {
		err = proto.Unmarshal(b, recv)
		if err != nil {
			t.logger.Critical(err.Error())
		}
		return
	}
	t.logger.Error(string(b))
	defer t.close()
}

func (t *transport) close() {
	err := t.tripper.Close()
	if err != nil {
		t.logger.Critical(err.Error())
	}
}
