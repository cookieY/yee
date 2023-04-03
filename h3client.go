package yee

import (
	"bytes"
	"crypto/tls"
	"github.com/quic-go/quic-go"
	"io/ioutil"
	"net/http"

	"github.com/cookieY/yee/logger"
	"github.com/golang/protobuf/proto"
	"github.com/quic-go/quic-go/http3"
)

type transport struct {
	addr               string
	insecureSkipVerify bool
	logger             logger.Logger
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
		logger:             logger.LogCreator(),
		c: &http.Client{
			Transport: tripper,
		},
		tripper: tripper,
	}
}

func (t *transport) Get(url string) (*http.Response, error) {
	resp, err := t.c.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (t *transport) Post(payload proto.Message, recv proto.Message) {
	p, err := proto.Marshal(payload)
	if err != nil {
		t.logger.Critical(err.Error())
		return
	}
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
