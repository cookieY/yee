package yee

import (
	"fmt"
	"net/http"
	"path"
)

type group struct {
	prefix     string
	middleware []HandlerFunc
	core       *Core
}

type UserFunc func(Context) error

func (g *group) Group(prefix string) *group {
	core := g.core
	newGroup := &group{
		prefix: g.prefix + prefix,
		core:   core,
	}
	g.core.groups = append(g.core.groups, newGroup)
	return newGroup
}

func (g *group) addRoute(method, prefix string, handler HandlerFunc) {
	pattern := g.prefix + prefix
	g.core.router.addRoute(method, pattern, handler)

}

func (g *group) Use(middleware ...HandlerFunc) {
	g.middleware = append(g.middleware, middleware...)
}

// todo: Implement the HTTP method and add router table

func (g *group) GET(path string, handler UserFunc) {
	g.addRoute("GET", path, HandlerFunc{Func: handler})
}

func (g *group) POST(path string, handler UserFunc) {
	g.addRoute("POST", path, HandlerFunc{Func: handler})
}

func (g *group) PUT(path string, handler UserFunc) {
	g.addRoute("PUT", path, HandlerFunc{Func: handler})
}

func (g *group) DELETE(path string, handler UserFunc) {
	g.addRoute("DELETE", path, HandlerFunc{Func: handler})
}

func (g *group) HEAD(path string, handler UserFunc) {
	g.addRoute("HEAD", path, HandlerFunc{Func: handler})
}

func (g *group) OPTION(path string, handler UserFunc) {
	g.addRoute("OPTION", path, HandlerFunc{Func: handler})
}

func (g *group) createDistHandler(relativePath string, fs http.FileSystem) UserFunc {
	ab := path.Join(g.prefix, relativePath)
	fs_1 := http.StripPrefix(ab, http.FileServer(fs))
	return func(c Context) (err error) {
		file := c.Params("filepath")
		if _, err := fs.Open(file); err != nil {
			return c.String(http.StatusNotFound, fmt.Sprintf("404 NOT FOUND!"))
		}
		fs_1.ServeHTTP(c.Response(), c.Request())
		return
	}

}

func (g *group) Static(relativePath, dist string) {
	handler := g.createDistHandler(relativePath, http.Dir(dist))
	url := path.Join(relativePath, "/*filepath")
	g.GET(url, handler)
}
