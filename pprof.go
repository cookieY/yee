package yee

import (
	"net/http"
	"net/http/pprof"
)

const DefaultPrefix = "/debug/pprof"

func getPrefix(prefixOptions string) string {
	prefix := DefaultPrefix
	if len(prefixOptions) > 1 {
		prefix += prefixOptions
	}
	return prefix
}

func WrapF(f http.HandlerFunc) HandlerFunc {
	return func(c Context) (err error) {
		f(c.Response(), c.Request())
		return nil
	}
}

func WrapH(h http.Handler) HandlerFunc {
	return func(c Context) (err error) {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func (c *Core) Pprof() {
	c.GET(getPrefix("/"), WrapF(pprof.Index))
	c.GET(getPrefix("/cmdline"), WrapF(pprof.Cmdline))
	c.GET(getPrefix("/profile"), WrapF(pprof.Profile))
	c.POST(getPrefix("/symbol"), WrapF(pprof.Symbol))
	c.GET(getPrefix("/symbol"), WrapF(pprof.Symbol))
	c.GET(getPrefix("/trace"), WrapF(pprof.Trace))
	c.GET(getPrefix("/allocs"), WrapH(pprof.Handler("allocs")))
	c.GET(getPrefix("/block"), WrapH(pprof.Handler("block")))
	c.GET(getPrefix("/goroutine"), WrapH(pprof.Handler("goroutine")))
	c.GET(getPrefix("/heap"), WrapH(pprof.Handler("heap")))
	c.GET(getPrefix("/mutex"), WrapH(pprof.Handler("mutex")))
	c.GET(getPrefix("/threadcreate"), WrapH(pprof.Handler("threadcreate")))
}
