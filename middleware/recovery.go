package middleware

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/cookieY/yee"
)

// Recovery is a recovery middleware
// when the program was panic
// it can recovery program and print stack info
func Recovery() yee.HandlerFunc {
	return func(c yee.Context) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				var pcs [32]uintptr
				n := runtime.Callers(3, pcs[:]) // skip first 3 caller

				var str strings.Builder
				str.WriteString("Traceback:")
				for _, pc := range pcs[:n] {
					fn := runtime.FuncForPC(pc)
					file, line := fn.FileLine(pc)
					str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
				}
				c.Logger().Critical(fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, str.String()))
				_ = c.ServerError(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		c.Next()
		return
	}
}
