package web

import (
	"html/template"
	"net/http"
	"path"
	"strings"
)

// HandlerFunc wrap handler function
// HandlerFunc version withou Context:
// `type HandlerFunc func(http.ResponseWriter, *http.Request)`
type HandlerFunc func(*Context)

// Engine implement Handler interface
type Engine struct {
	// router *router

	// RouterGroup pointer to access by engine
	*RouterGroup
	groups []*RouterGroup

	// Template Render
	template *template.Template
	funcMap  template.FuncMap
}

// * Version Without RouterGroup
// NewEngine create Engine instance
/*
func NewEngine() *Engine {
	return &Engine{
		router: newRouter(),
	}
}
*/

// NewEngine create Engine instance
func NewEngine() *Engine {
	engine := &Engine{RouterGroup: &RouterGroup{router: newRouter()}}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Run start the http server with engine
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

// Router return engine's router info
// TODO: return router map info
func (e *Engine) Router() *router {
	return e.router
}

func (e *Engine) addRouter(method, path string, handler HandlerFunc) {
	// key := method + "-" + path
	// e.router[key] = handler
	e.router.addRoute(method, path, handler)
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
	// * Version: With Context
	ctx := newContext(w, req)
	ctx.handlers = e.getMiddlewares(w, req)
	ctx.engine = e
	e.router.handle(ctx)
}

// getMiddlewares return all router groups middlewares
func (e *Engine) getMiddlewares(w http.ResponseWriter, req *http.Request) []HandlerFunc {
	middlewares := []HandlerFunc{}
	for _, g := range e.groups {
		if strings.HasPrefix(req.URL.Path, g.prefix) {
			middlewares = append(middlewares, g.middlewares...)
		}
	}

	return middlewares
}

// SetFuncMap setting funcMap for engine
func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

// LoadGlobHTML for global html template
func (e *Engine) LoadGlobHTML(pattern string) {
	e.template = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}

// RouterGroup group router with prefix
type RouterGroup struct {
	prefix      string
	router      *router
	middlewares []HandlerFunc
}

// Group create new RouterGroup instance with prefix.
// Every RouterGroup instance share the same Engine.
func (g *RouterGroup) Group(prefix string) *RouterGroup {
	return &RouterGroup{
		prefix: g.prefix + prefix,
		router: g.router,
	}
}

func (g *RouterGroup) addRoute(method, prefix string, handler HandlerFunc) {
	pattern := g.prefix + prefix
	g.router.addRoute(method, pattern, handler)
}

// GET handler http get request by group
func (g *RouterGroup) GET(path string, handler HandlerFunc) {
	g.addRoute(http.MethodGet, path, handler)
}

// POST handler http post request by group
func (g *RouterGroup) POST(path string, handler HandlerFunc) {
	g.addRoute(http.MethodPost, path, handler)
}

// PUT handler http put request by group
func (g *RouterGroup) PUT(path string, handler HandlerFunc) {
	g.addRoute(http.MethodPut, path, handler)
}

// DELETE handler http delete request by group
func (g *RouterGroup) DELETE(path string, handler HandlerFunc) {
	g.addRoute(http.MethodDelete, path, handler)
}

// Use register serial middlewares to the router group
func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// createStaticHandler for making static handler
func (g *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(g.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(ctx *Context) {
		file := ctx.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.StatusCode = http.StatusOK
		fileServer.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// Static serve static files
func (g *RouterGroup) Static(relativePath, root string) {
	handler := g.createStaticHandler(relativePath, http.Dir(root))
	pattern := path.Join(relativePath, "/*filepath")
	g.GET(pattern, handler)
}
