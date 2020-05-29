package yee

import (
	"net/http"
	"sync"
)

type HandlerFunc func(Context) error

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
	Logger                 *Logger
}

type HTTPError struct {
	Code     int
	Message  interface{}
	Internal error // Stores the error returned by an external dependency
}

const YeeVersion = "Yee v0.0.1"

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
		Logger: &Logger{},
	}
	core.core = core
	core.pool.New = func() interface{} {
		return core.allocateContext()
	}
	return core
}
func (c *Core) allocateContext() *context {
	v := make(Params, 0, c.maxParams)
	return &context{engine: c, params: &v}
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
			//c.writermem.WriteHeaderNow()
			return
		}
		//if httpMethod != "CONNECT" && rPath != "/" {
		//	if value.tsr && engine.RedirectTrailingSlash {
		//		redirectTrailingSlash(c)
		//		return
		//	}
		//	if engine.RedirectFixedPath && redirectFixedPath(c, root, engine.RedirectFixedPath) {
		//		return
		//	}
		//}
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
	_ = context.String(404, "NOT FOUND")
}

//func serveError(c *context, code int, defaultMessage []byte) {
//	c.writermem.status = code
//	c.Next()
//	if c.writermem.Written() {
//		return
//	}
//	if c.writermem.Status() == code {
//		c.writermem.Header()["Content-Type"] = mimePlain
//		_, err := c.Writer.Write(defaultMessage)
//		if err != nil {
//			debugPrint("cannot write message to writer during serve error: %v", err)
//		}
//		return
//	}
//	c.writermem.WriteHeaderNow()
//}
