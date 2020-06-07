package yee

import (
	"net/http"
	"path"
	"sync"
)

type HandlerFunc struct {
	Func         UserFunc
	IsMiddleware bool
}

type UserFunc func(Context) error

type HandlersChain []HandlerFunc

// Core implement  httpServer interface
type Core struct {
	*router
	trees                  methodTrees
	pool                   sync.Pool
	maxParams              uint16
	HandleMethodNotAllowed bool
	allNoRoute             HandlersChain
	allNoMethod            HandlersChain
	noRoute                HandlersChain
	l                      logger
	RedirectTrailingSlash  bool
	RedirectFixedPath      bool
}

type HTTPError struct {
	Code     int
	Message  interface{}
	Internal error // Stores the error returned by an external dependency
}

const Version = "Yee v0.0.1"

// init Core
func New() *Core {
	router := &router{
		handlers: nil,
		root:     true,
		basePath: "/",
	}

	core := &Core{
		trees:  make(methodTrees, 0, 0),
		router: router,
		l:      logger{level: 6},
	}
	core.core = core
	core.pool.New = func() interface{} {
		return core.allocateContext()
	}
	return core
}
func (c *Core) allocateContext() *context {
	v := make(Params, 0, c.maxParams)
	return &context{engine: c, params: &v, index: -1}
}

func (c *Core) Use(middleware ...HandlerFunc) {
	c.router.Use(middleware...)
}

// override Handler.ServeHTTP
func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := c.pool.Get().(*context)
	context.writermem.reset(w)
	context.r = r
	context.reset()

	c.handleHTTPRequest(context)

	c.pool.Put(context)
}

func (c *Core) Start(addr string) {
	if err := http.ListenAndServe(addr, c); err != nil {
		panic(err)
	}
}

// NewHTTPError creates a new HTTPError instance.
func NewHTTPError(code int, message ...interface{}) *HTTPError {
	he := &HTTPError{Code: code, Message: http.StatusText(code)}
	if len(message) > 0 {
		he.Message = message[0]
	}
	return he
}

func (c *Core) handleHTTPRequest(context *context) {
	httpMethod := context.r.Method
	rPath := context.r.URL.Path
	unescape := false
	//if engine.UseRawPath && len(context.Request().URL.RawPath) > 0 {
	//	rPath = c.Request.URL.RawPath
	//	unescape = engine.UnescapePathValues
	//}
	//
	//if engine.RemoveExtraSlash {
	//	rPath = cleanPath(rPath)
	//}

	// Find root of the tree for the given HTTP method
	t := c.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method != httpMethod {
			continue
		}
		root := t[i].root
		// Find route in tree
		value := root.getValue(rPath, context.params, unescape)
		if value.params != nil {
			context.Param = *value.params
		}
		if value.handlers != nil {
			context.handlers = value.handlers
			context.path = value.fullPath
			context.Next()
			context.writermem.WriteHeaderNow()
			return
		}
		if httpMethod != "CONNECT" && rPath != "/" {
			if value.tsr && c.RedirectTrailingSlash {
				redirectTrailingSlash(context)
				return
			}
			//if c.RedirectFixedPath && redirectFixedPath(c, root, c.RedirectFixedPath) {
			//	return
			//}
		}
		break
	}

	//if c.HandleMethodNotAllowed {
	//	for _, tree := range c.trees {
	//		if tree.method == httpMethod {
	//			continue
	//		}
	//		if value := tree.root.getValue(rPath, nil, unescape); value.handlers != nil {
	//			c.handlers = c.allNoMethod
	//			serveError(c, http.StatusMethodNotAllowed, default405Body)
	//			return
	//		}
	//	}
	//}
	context.handlers = c.allNoRoute
	context.ServerError(404, []byte("404 NOT FOUND"),false)
}

func redirectTrailingSlash(c *context) {
	req := c.r
	p := req.URL.Path
	if prefix := path.Clean(c.r.Header.Get("X-Forwarded-Prefix")); prefix != "." {
		p = prefix + "/" + req.URL.Path
	}
	req.URL.Path = p + "/"
	if length := len(p); length > 1 && p[length-1] == '/' {
		req.URL.Path = p[:length-1]
	}
	redirectRequest(c)
}

//func redirectFixedPath(c *context, root *node, trailingSlash bool) bool {
//	req := c.r
//	rPath := req.URL.Path
//
//	if fixedPath, ok := root.findCaseInsensitivePath(cleanPath(rPath), trailingSlash); ok {
//		req.URL.Path = BytesToString(fixedPath)
//		redirectRequest(c)
//		return true
//	}
//	return false
//}

func redirectRequest(c *context) {
	req := c.r
	rURL := req.URL.String()

	code := http.StatusMovedPermanently // Permanent redirect, request with GET method
	if req.Method != http.MethodGet {
		code = http.StatusTemporaryRedirect
	}
	http.Redirect(c.w, req, rURL, code)
	c.writermem.WriteHeaderNow()
}