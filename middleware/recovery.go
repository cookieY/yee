package middleware

import (
	"fmt"
	"github.com/cookieY/Yee"
	"net/http"
	"runtime"
	"strings"
)

func trace(info string) string {
	var pcs [12]uintptr
	n := runtime.Callers(3, pcs[:])
	var str strings.Builder
	str.WriteString(info + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() yee.HandlerFunc {
	return yee.HandlerFunc{
		Func: func(c yee.Context) (err error) {
			defer func() {
				if err := recover(); err != nil {
					message := fmt.Sprintf("%s", err)
					c.Logger().Critical(trace(message))
					c.ServerError(http.StatusInternalServerError, []byte("Internal Server Error"), true)
				}
			}()
			c.Next()
			return
		},
		IsMiddleware: true,
	}
}
