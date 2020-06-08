package middleware

import (
	"fmt"
	"github.com/cookieY/yee"
	"github.com/valyala/fasttemplate"
	"io"
	"log"
)

type (
	LoggerConfig struct {
		Format   string
		Level    uint8
		IsLogger bool
	}
)

var DefaultLoggerConfig = LoggerConfig{
	Format:   `"url":"${url}" "method":"${method}" "status":${status} "protocol":"${protocol}" "remote_ip":"${remote_ip}" "bytes_in": "${bytes_in} bytes" "bytes_out": "${bytes_out} bytes"`,
	Level:    3,
	IsLogger: true,
}

func Logger() yee.HandlerFunc {
	return LoggerWithConfig(DefaultLoggerConfig)
}

func LoggerWithConfig(config LoggerConfig) yee.HandlerFunc {
	if config.Format == "" {
		config.Format = DefaultLoggerConfig.Format
	}

	if config.Level == 0 {
		config.Level = DefaultLoggerConfig.Level
	}

	t, err := fasttemplate.NewTemplate(config.Format, "${", "}")

	if err != nil {
		log.Fatalf("unexpected error when parsing template: %s", err)
	}

	logger := yee.LogCreator()

	logger.SetLevel(config.Level)

	logger.IsLogger(config.IsLogger)

	return yee.HandlerFunc{
		Func: func(context yee.Context) (err error) {
			context.Next()
			s := t.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
				switch tag {
				case "url":
					p := context.Request().URL.Path
					if p == "" {
						p = "/"
					}
					return w.Write([]byte(p))
				case "method":
					return w.Write([]byte(context.Request().Method))
				case "status":
					return w.Write([]byte(fmt.Sprintf("%d", context.Response().Status())))
				case "remote_ip":
					return w.Write([]byte(context.RemoteIp()))
				case "host":
					return w.Write([]byte(context.Request().Host))
				case "protocol":
					return w.Write([]byte(context.Request().Proto))
				case "bytes_in":
					cl := context.Request().Header.Get(yee.HeaderContentLength)
					if cl == "" {
						cl = "0"
					}
					return w.Write([]byte(cl))
				case "bytes_out":
					return w.Write([]byte(fmt.Sprintf("%d", context.Response().Size())))
				default:
					return w.Write([]byte(""))
				}
			})
			if context.Response().Status() < 400 {
				logger.Trace(s)
			} else {
				logger.Warn(s)
			}
			return
		},
		IsMiddleware: true,
	}
}
