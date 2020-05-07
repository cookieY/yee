package knocker

import (
	"net/http"
	"strings"
)

type HandlerFunc func(Context) error

// Core implement  httpServer interface
type Core struct {
	router *router
	*group
	groups []*group
}

type HTTPError struct {
	Code     int
	Message  interface{}
	Internal error // Stores the error returned by an external dependency
}


// init Core
func New() *Core {
	core := &Core{router: newRouter()}
	core.group = &group{core: core}
	core.groups = []*group{core.group}
	return core
}

// override Handler.ServeHTTP
func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandlerFunc
	for _, g := range c.groups {
		if strings.HasPrefix(r.URL.Path, g.prefix) {
			middlewares = append(middlewares, g.middleware...)
		}
	}
	context := newContext(w, r)
	context.handlers = middlewares
	c.router.handle(context)
}

func (c *Core) Start(addr string) error {
	return http.ListenAndServe(addr, c)
}

// NewHTTPError creates a new HTTPError instance.
func NewHTTPError(code int, message ...interface{}) *HTTPError {
	he := &HTTPError{Code: code, Message: http.StatusText(code)}
	if len(message) > 0 {
		he.Message = message[0]
	}
	return he
}
