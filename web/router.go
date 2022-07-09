package web

import (
	"log"
	"net/http"
	"strings"
)

// * Version of Static Match Router
/*
type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

func (r *router) addRouter(method, path string, handler HandlerFunc) {
	log.Printf("[Router] %4s - %s", method, path)
	key := method + "-" + path
	r.handlers[key] = handler
}

func (r *router) handle(ctx *Context) {
	key := ctx.Method + "-" + ctx.Path
	if handler, ok := r.handlers[key]; ok {
		// handler(ctx.Writer, ctx.Request)
		handler(ctx)
	} else {
		ctx.String(http.StatusNotFound, "404 Not Found: %q\n", ctx.Path)
	}
}
*/

// * Version of Trie Tree Match Router
// router storage roots and handlers info
// roots key e.g: roots["GET"] root["POST"]
// handlers e.g: handlers["GET-/p/:lang/doc"]
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	subs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range subs {
		if item != "" {
			parts = append(parts, item)
			if string(item[0]) == "*" {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method, path string, handler HandlerFunc) {
	parts := parsePattern(path)
	key := method + "-" + path
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(path, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	log.Printf("searchParts: %v", searchParts)
	params := make(map[string]string)
	root, ok := r.roots[method]
	log.Printf("root: %v, ok: %v", root, ok)
	if !ok {
		return nil, nil
	}

	node := root.search(searchParts, 0)
	log.Printf("searched node: %v", node)
	if node != nil {
		parts := parsePattern(node.pattern)
		log.Printf("searched parts: %v", parts)
		for i, part := range parts {
			if string(part[0]) == ":" {
				params[part[1:]] = searchParts[i]
			}
			if string(part[0]) == "*" && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[i:], "/")
				break
			}
		}
		return node, params
	}

	return nil, nil
}

func (r *router) handle(ctx *Context) {
	node, params := r.getRoute(ctx.Method, ctx.Path)
	if node != nil {
		ctx.Params = params
		key := ctx.Method + "-" + node.pattern
		r.handlers[key](ctx)
	} else {
		ctx.String(http.StatusNotFound, "404 Not Found: %q\n", ctx.Path)
	}
}
