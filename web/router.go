package web

import (
	"log"
	"net/http"
)

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
