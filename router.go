package yee

import (
	"net/http"
	"path"
	"strings"
)

type Router struct {
	handlers []HandlerFunc
	core     *Core
	root     bool
	basePath string
}

// RestfulAPI is the default implementation of restfulApi interface
type RestfulAPI struct {
	Get    HandlerFunc
	Post   HandlerFunc
	Delete HandlerFunc
	Put    HandlerFunc
}

// Implement the HTTP method and add to the router table
// GET,POST,PUT,DELETE,OPTIONS,TRACE,HEAD,PATCH
// these are defined in RFC 7231 section 4.3.

func (r *Router) GET(path string, handler ...HandlerFunc) {
	r.handle(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler ...HandlerFunc) {
	r.handle(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler ...HandlerFunc) {
	r.handle(http.MethodPut, path, handler)
}

func (r *Router) DELETE(path string, handler ...HandlerFunc) {
	r.handle(http.MethodDelete, path, handler)
}

func (r *Router) PATCH(path string, handler ...HandlerFunc) {
	r.handle(http.MethodPatch, path, handler)
}

func (r *Router) HEAD(path string, handler ...HandlerFunc) {
	r.handle(http.MethodHead, path, handler)
}

func (r *Router) TRACE(path string, handler ...HandlerFunc) {
	r.handle(http.MethodTrace, path, handler)
}

func (r *Router) OPTIONS(path string, handler ...HandlerFunc) {
	r.handle(http.MethodOptions, path, handler)
}

func (r *Router) Restful(path string, api RestfulAPI) {

	if api.Get != nil {
		r.handle(http.MethodGet, path, HandlersChain{api.Get})
	}
	if api.Post != nil {
		r.handle(http.MethodPost, path, HandlersChain{api.Post})
	}
	if api.Put != nil {
		r.handle(http.MethodPut, path, HandlersChain{api.Put})
	}
	if api.Delete != nil {
		r.handle(http.MethodDelete, path, HandlersChain{api.Delete})
	}
}

func (r *Router) Any(path string, handler ...HandlerFunc) {
	r.handle(http.MethodPost, path, handler)
	r.handle(http.MethodGet, path, handler)
	r.handle(http.MethodPut, path, handler)
	r.handle(http.MethodDelete, path, handler)
	r.handle(http.MethodOptions, path, handler)
}

func (r *Router) Use(middleware ...HandlerFunc) {
	r.handlers = append(r.handlers, middleware...)
}

func (r *Router) Group(prefix string, handlers ...HandlerFunc) *Router {
	rx := &Router{
		handlers: r.combineHandlers(handlers),
		core:     r.core,
		basePath: r.calculateAbsolutePath(prefix),
	}
	return rx
}

func (r *Router) handle(method, path string, handlers HandlersChain) {
	absolutePath := r.calculateAbsolutePath(path)
	handlers = r.combineHandlers(handlers)
	r.core.addRoute(method, absolutePath, handlers)
}

func (c *Core) addRoute(method, prefix string, handlers HandlersChain) {
	if prefix[0] != '/' {
		panic("path must begin with '/'")
	}

	if method == "" {
		panic("HTTP method can not be empty")
	}

	root := c.trees.get(method)
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		c.trees = append(c.trees, methodTree{method: method, root: root})
	}
	root.addRoute(prefix, handlers)

	// Update maxParams
	if paramsCount := countParams(prefix); paramsCount > c.maxParams {
		c.maxParams = paramsCount
	}

}

func (r *Router) Static(relativePath, root string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL path cannot be used when serving a static folder")
	}
	handler := r.createDistHandler(relativePath, http.Dir(root))
	url := path.Join(relativePath, "/*filepath")
	r.GET(url, handler)
	r.HEAD(url, handler)

}

func (r *Router) Packr(relativePath string, fs http.FileSystem) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL path cannot be used when serving a static folder")
	}
	handler := r.createDistHandler(relativePath, fs)
	url := path.Join(relativePath, "/*filepath")
	r.GET(url, handler)
	r.HEAD(url, handler)

}

func (r *Router) createDistHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := r.calculateAbsolutePath(relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c Context) (err error) {
		file := c.Params("filepath")
		// Mmap provides a way to memory-map,
		// This improves reading efficiency when a large number of static file reads are performed
		f, err2 := fs.Open(file)
		if err2 != nil {
			c.Status(http.StatusNotFound)
			c.Reset()
		}
		if f != nil {
			_ = f.Close()
		}
		fileServer.ServeHTTP(c.Response(), c.Request())
		return
	}
}

func (r *Router) calculateAbsolutePath(relativePath string) string {
	return joinPaths(r.basePath, relativePath)
}

func (r *Router) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(r.handlers) + len(handlers)
	if finalSize >= int(abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, r.handlers)
	copy(mergedHandlers[len(r.handlers):], handlers)
	return mergedHandlers
}
