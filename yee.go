package yee

import (
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net/http"
	"os"
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
	H2server               *http.Server
	allNoRoute             HandlersChain
	allNoMethod            HandlersChain
	noRoute                HandlersChain
	noMethod               HandlersChain
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

const banner = `
  ___    ___ _______   _______      
 |\  \  /  /|\  ___ \ |\  ___ \     
 \ \  \/  / | \   __/|\ \   __/|    
  \ \    / / \ \  \_|/_\ \  \_|/__  
   \/  /  /   \ \  \_|\ \ \  \_|\ \ 
 __/  / /      \ \_______\ \_______\
|\___/ /        \|_______|\|_______|  %s
\|___|/

`

// init Core
func InitCore() *Core {
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

func New() *Core {
	fmt.Printf(banner, Version)
	return InitCore()
}

func (c *Core) allocateContext() *context {
	v := make(Params, 0, c.maxParams)
	return &context{engine: c, params: &v, index: -1}
}

func (c *Core) Use(middleware ...HandlerFunc) {
	c.router.Use(middleware...)
	c.rebuild404Handlers()
	c.rebuild405Handlers()

}

func (c *Core) rebuild404Handlers() {
	c.allNoRoute = c.combineHandlers(c.noRoute)
}

func (c *Core) rebuild405Handlers() {
	c.allNoMethod = c.combineHandlers(c.noMethod)
}

// override Handler.ServeHTTP
// all requests/response deal with here
// we use sync.pool save context variable
// because we do this can be used less memory
// we just only reset context, when before callback c.handleHTTPRequest func
// and put context variable into poll

func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := c.pool.Get().(*context)
	context.writermem.reset(w)
	context.r = r
	context.reset()

	c.handleHTTPRequest(context)

	c.pool.Put(context)
}

func (c *Core) Run(addr string) {
	if err := http.ListenAndServe(addr, c); err != nil {
		c.l.Critical(err.Error())
		os.Exit(1)
	}
}


// golang supports http2,if client supports http2
// Otherwise, the http protocol return to http1.1
func (c *Core) RunTLS(addr, certFile, keyFile string) {
	if err := http.ListenAndServeTLS(addr, certFile, keyFile, c); err != nil {
		c.l.Critical(err.Error())
		os.Exit(1)
	}
}

// In normal conditions, http2 must used certificate
// H2C is non-certificate`s http2
// notes:
// 1.the browser is not supports H2C proto, you should write your web client test program
// 2.the H2C protocol is not safety
func (c *Core) RunH2C(addr string) {
	s := &http2.Server{}
	h1s := &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(c, s),
	}
	log.Fatal(h1s.ListenAndServe())
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

		break
	}

	context.handlers = c.allNoRoute

	// Notice
	// We must judge whether an empty request is OPTIONS method,
	// Because when complex request (XMLHttpRequest) will send an OPTIONS request and fetch the preflight resource.
	// But in general, we do not register an OPTIONS handle,
	// So this may cause some middleware errors.

	if httpMethod == http.MethodOptions {
		serveError(context, 200, []byte("preflight resource"))
	} else {
		serveError(context, http.StatusNotFound, []byte("404 NOT FOUND"))
	}
}

func serveError(c *context, code int, defaultMessage []byte) {
	c.writermem.status = code
	c.Next()
	if c.writermem.Written() {
		return
	}
	if c.writermem.Status() == code {
		c.writermem.Header()["Content-Type"] = []string{MIMETextPlain}
		_, err := c.w.Write(defaultMessage)
		if err != nil {
			c.engine.l.Error(fmt.Sprintf("cannot write message to writer during serve error: %v", err))
		}
		return
	}
	c.writermem.WriteHeaderNow()
}