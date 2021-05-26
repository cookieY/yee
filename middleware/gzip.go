package middleware

import (
	"bufio"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/cookieY/yee"
)

// GzipConfig defines config of Gzip middleware
type GzipConfig struct {
	Level int
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// DefaultGzipConfig is the default config of gzip middleware
var DefaultGzipConfig = GzipConfig{Level: 1}

// Gzip is the default implementation of gzip middleware
func Gzip() yee.HandlerFunc {
	return GzipWithConfig(DefaultGzipConfig)
}

// GzipWithConfig is the custom implementation of gzip middleware
func GzipWithConfig(config GzipConfig) yee.HandlerFunc {
	if config.Level == 0 {
		config.Level = DefaultGzipConfig.Level
	}

	return func(c yee.Context) (err error) {
		if c.IsWebsocket() {
			return
		}
		res := c.Response()
		res.Header().Add(yee.HeaderVary, yee.HeaderAcceptEncoding)
		if strings.Contains(c.Request().Header.Get(yee.HeaderAcceptEncoding), "gzip") {
			res.Header().Set(yee.HeaderContentEncoding, "gzip")
			rw := res.Writer()
			w, err := gzip.NewWriterLevel(rw, config.Level)
			if err != nil {
				return err
			}
			defer func() {
				if res.Size() < 1 {
					if res.Header().Get(yee.HeaderContentEncoding) == "gzip" {
						res.Header().Del(yee.HeaderContentEncoding)
					}
					res.Override(rw)
					w.Reset(ioutil.Discard)
				}
				_ = w.Close()
			}()
			grw := &gzipResponseWriter{Writer: w, ResponseWriter: rw}
			res.Override(grw)
		}
		c.Next()
		return
	}
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if w.Header().Get(yee.HeaderContentType) == "" {
		w.Header().Set(yee.HeaderContentType, http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}

func (w *gzipResponseWriter) Flush() {
	w.Writer.(*gzip.Writer).Flush()
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}
