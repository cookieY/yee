package middleware

import (
	"bufio"
	"compress/gzip"
	"github.com/cookieY/yee"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

type GzipConfig struct {
	Level int
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

var DefaultGzipConfig = GzipConfig{Level: 1}

func Gzip() yee.HandlerFunc {
	return GzipWithConfig(DefaultGzipConfig)
}

func GzipWithConfig(config GzipConfig) yee.HandlerFunc {
	if config.Level == 0 {
		config.Level = DefaultGzipConfig.Level
	}

	return yee.HandlerFunc{
		Func: func(c yee.Context) (err error) {
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
		},
		IsMiddleware: true,
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
