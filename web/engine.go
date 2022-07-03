package web

import (
	"net/http"
)

// HandlerFunc wrap handler function
// HandlerFunc version withou Context:
// `type HandlerFunc func(http.ResponseWriter, *http.Request)`
type HandlerFunc func(*Context)

// Engine implement Handler interface
type Engine struct {
	// router map[string]HandlerFunc
	router *router
}

// NewEngine create Engine instance
func NewEngine() *Engine {
	return &Engine{
		// router: make(map[string]HandlerFunc),
		router: newRouter(),
	}
}

// Run start the http server with engine
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) addRouter(method, path string, handler HandlerFunc) {
	// key := method + "-" + path
	// e.router[key] = handler
	e.router.addRouter(method, path, handler)
}

// GET handler http get request
func (e *Engine) GET(path string, handler HandlerFunc) {
	e.addRouter(http.MethodGet, path, handler)
}

// POST handler http post request
func (e *Engine) POST(path string, handler HandlerFunc) {
	e.addRouter(http.MethodPost, path, handler)
}

// PUT handler http put request
func (e *Engine) PUT(path string, handler HandlerFunc) {
	e.addRouter(http.MethodPut, path, handler)
}

// DELETE handler http delete request
func (e *Engine) DELETE(path string, handler HandlerFunc) {
	e.addRouter(http.MethodDelete, path, handler)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// * Version: Without Context
	// key := req.Method + "-" + req.URL.Path
	// if handler, ok := e.router[key]; ok {
	// 	handler(w, req)
	// } else {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	fmt.Fprintf(w, "404 Not FOUND: %s\n", req.URL)
	// }

	// * Version: With Context
	ctx := newContext(w, req)
	e.router.handle(ctx)
}
