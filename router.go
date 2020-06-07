package yee

import (
	"net/http"
	"path"
	"strings"
)

type router struct {
	handlers []HandlerFunc
	core     *Core
	root     bool
	basePath string
}

// todo: Implement the HTTP method and add router table

func (r *router) GET(path string, handler UserFunc) {
	r.handle(http.MethodGet, path, []HandlerFunc{{Func: handler}})
}

func (r *router) POST(path string, handler UserFunc) {
	r.handle(http.MethodPost, path, []HandlerFunc{{Func: handler}})
}

func (r *router) PATCH(path string, handler UserFunc) {
	r.handle(http.MethodPatch, path, []HandlerFunc{{Func: handler}})
}

func (r *router) PUT(path string, handler UserFunc) {
	r.handle(http.MethodPut, path, []HandlerFunc{{Func: handler}})
}

func (r *router) DELETE(path string, handler UserFunc) {
	r.handle(http.MethodDelete, path, []HandlerFunc{{Func: handler}})
}

func (r *router) HEAD(path string, handler UserFunc) {
	r.handle(http.MethodHead, path, []HandlerFunc{{Func: handler}})
}

func (r *router) OPTIONS(path string, handler UserFunc) {
	r.handle(http.MethodOptions, path, []HandlerFunc{{Func: handler}})
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

func (r *router) Static(relativePath, root string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL path cannot be used when serving a static folder")
	}
	handler := r.createDistHandler(relativePath, http.Dir(root))
	url := path.Join(relativePath, "/*filepath")
	r.GET(url, handler)
	r.HEAD(url, handler)

}

func (r *router) createDistHandler(relativePath string, fs http.FileSystem) UserFunc {
	absolutePath := r.calculateAbsolutePath(relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c Context) (err error) {
		file := c.Params("filepath")
		f, err2 := fs.Open(file)
		if err2 != nil {
			c.Status(http.StatusNotFound)
			c.Reset()
		} else {
			_ = f.Close()
			fileServer.ServeHTTP(c.Response(), c.Request())
		}
		return
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
