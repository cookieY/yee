package yee

import (
	"net/http"
)

type router struct {
	handlers []HandlerFunc
	core     *Core
	root     bool
	basePath string
}

// todo: Implement the HTTP method and add router table

func (r *router) GET(path string, handler ...HandlerFunc) {
	r.handle(http.MethodGet, path, handler)
}

func (r *router) POST(path string, handler ...HandlerFunc) {
	r.handle(http.MethodPost, path, handler)
}

func (r *router) PATCH(path string, handler ...HandlerFunc) {
	r.handle(http.MethodPatch, path, handler)
}

func (r *router) PUT(path string, handler ...HandlerFunc) {
	r.handle(http.MethodPut, path, handler)
}

func (r *router) DELETE(path string, handler ...HandlerFunc) {
	r.handle(http.MethodDelete, path, handler)
}

func (r *router) HEAD(path string, handler ...HandlerFunc) {
	r.handle(http.MethodHead, path, handler)
}

func (r *router) OPTIONS(path string, handler ...HandlerFunc) {
	r.handle(http.MethodOptions, path, handler)
}

func (r *router) Use(middleware ...HandlerFunc) {
	r.handlers = append(r.handlers, middleware...)
}

func (r *router) Group(prefix string, handlers ...HandlerFunc) *router {
	return &router{
		handlers: r.combineHandlers(handlers),
		core:     r.core,
		basePath: r.calculateAbsolutePath(prefix),
	}
}

func (r *router) handle(method, path string, handlers HandlersChain) {
	absolutePath := r.calculateAbsolutePath(path)
	handlers = r.combineHandlers(handlers)
	r.core.addRoute(method, absolutePath, handlers)
}

func (c *Core) addRoute(method, prefix string, handlers HandlersChain) {
	assertS(prefix[0] == '/', "path must begin with '/'")
	assertS(method != "", "HTTP method can not be empty")
	assertS(len(handlers) > 0, "there must be at least one handler")

	//debugPrintRoute(method, path, handlers)

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

func (r *router) HTTPHandlerFunc(h http.HandlerFunc) HandlerFunc {
	return func(c Context) error {
		h(c.Response(), c.Request())
		return nil
	}
}

func (r *router) calculateAbsolutePath(relativePath string) string {
	return joinPaths(r.basePath, relativePath)
}

func (r *router) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(r.handlers) + len(handlers)
	if finalSize >= int(abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, r.handlers)
	copy(mergedHandlers[len(r.handlers):], handlers)
	return mergedHandlers
}
