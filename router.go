package yee

import (
	"fmt"
	"net/http"
	"strings"
)

type router struct {
	route    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
		route:    make(map[string]*node),
	}
}

func parserParts(pattern string) []string {
	l := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, i := range l {
		if i != "" {
			parts = append(parts, i)
			if i[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method, path string, handler HandlerFunc) {

	parts := parserParts(path)

	handlePath := fmt.Sprintf("%s-%s", method, path)

	if _, ok := r.route[method]; !ok {
		r.route[method] = &node{}
	}
	r.route[method].insert(path, parts, 0)

	r.handlers[handlePath] = handler
}

func (r *router) fetchRoute(method, path string) (*node, map[string]string) {
	sParts := parserParts(path)
	params := make(map[string]string)
	root, ok := r.route[method]
	if !ok {
		return nil, nil
	}

	n := root.search(sParts, 0) // 查找子节点node列表 从cn开始
	if n != nil {               // 如果存在节点返回节点信息以及params
		parts := parserParts(n.pattern)
		for idx, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = sParts[idx]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(sParts[idx:], "/")
				break
			}
			if part[0] == '/' && len(part) > 1 {
				params["*"] = strings.Join(sParts[idx:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// Process all requests and into the router table
func (r *router) handle(context *context) {
	n, params := r.fetchRoute(context.method, context.path)
	if n != nil {
		context.params = params
		path := fmt.Sprintf("%s-%s", context.method, n.pattern)
		context.handlers = append(context.handlers, r.handlers[path])
	} else {
		context.handlers = append(context.handlers, HandlerFunc{
			Func: func(c Context) (err error) {
				return c.String(http.StatusNotFound, fmt.Sprintf("404 NOT FOUND: %s\n", context.path))
			},
			IsMiddleware: false,
		})
	}
	context.Next()
}
