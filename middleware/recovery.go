package middleware
//
//import (
//	"fmt"
//	"log"
//	"net/http"
//	"runtime"
//	"strings"
//	"yee"
//)
//
//func trace(info string) string {
//	var pcs [12]uintptr
//	n := runtime.Callers(5, pcs[:])
//	var str strings.Builder
//	str.WriteString(info + "\nTraceback:")
//	for _, pc := range pcs[:n] {
//		fn := runtime.FuncForPC(pc)
//		file, line := fn.FileLine(pc)
//		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
//	}
//	return str.String()
//}
//
//func Recovery() yee.HandlerFunc {
//	return func(c yee.Context) (err error) {
//		defer func() {
//			if err := recover(); err != nil {
//				message := fmt.Sprintf("%s", err)
//				log.Printf("%s", trace(message))
//				_ = c.String(http.StatusInternalServerError, "Internal Server Error")
//			}
//		}()
//		c.Next()
//		return
//	}
//}
