package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H alias for map
type H map[string]interface{}

// Context of request
type Context struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	Path       string
	Method     string
	StatusCode int
	Params     map[string]string

	// Middleware Properties
	handlers []HandlerFunc
	index    int

	// Global Engine Instance Pointer
	engine *Engine
}

// newContext contains http request and response info
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Request: req,
		Writer:  w,
		Method:  req.Method,
		Path:    req.URL.Path,
		index:   -1,
	}
}

// Param return route param by name
func (ctx *Context) Param(key string) string {
	value, _ := ctx.Params[key]
	return value
}

// PostForm return post body value
func (ctx *Context) PostForm(key string) string {
	return ctx.Request.FormValue(key)
}

// Query return request raw query value
func (ctx *Context) Query(key string) string {
	return ctx.Request.URL.Query().Get(key)
}

// Status write response header with status code
func (ctx *Context) Status(code int) {
	ctx.StatusCode = code
	ctx.Writer.WriteHeader(code)
}

// SetHeader sets header entries of response
func (ctx *Context) SetHeader(key, value string) {
	ctx.Writer.Header().Set(key, value)
}

// String return a raw string response
func (ctx *Context) String(code int, format string, values ...interface{}) {
	ctx.Status(code)
	ctx.SetHeader("Content-Type", "text/plain; charset=utf-8")
	ctx.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// Data return a raw string response with bytes param
func (ctx *Context) Data(code int, value []byte) {
	ctx.Status(code)
	ctx.SetHeader("Content-Type", "text/plain; charset=utf-8")
	ctx.Writer.Write(value)
}

// JSON return a json string response
func (ctx *Context) JSON(code int, value interface{}) {
	ctx.Status(code)
	ctx.SetHeader("Content-Type", "application/json")
	enc := json.NewEncoder(ctx.Writer)
	if err := enc.Encode(value); err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// HTML return a html string response
// * Version without html support
/*
func (ctx *Context) HTML(code int, html string) {
	ctx.Status(code)
	ctx.SetHeader("Content-Type", "text/html")
	ctx.Writer.Write([]byte(html))
}
*/

// HTML return a html string response
func (ctx *Context) HTML(code int, name string, data interface{}) {
	ctx.Status(code)
	ctx.SetHeader("Content-Type", "text/html")
	err := ctx.engine.template.ExecuteTemplate(ctx.Writer, name, data)
	if err != nil {
		// ctx.String(http.StatusInternalServerError, "%v", err)
		ctx.Fail(http.StatusInternalServerError, err.Error())
	}
}

// Next process middleware handlers serially
func (ctx *Context) Next() {
	ctx.index++
	l := len(ctx.handlers)
	for ; ctx.index < l; ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}

// Fail handle error response
func (ctx *Context) Fail(code int, err string) {
	ctx.index = len(ctx.handlers)
	ctx.JSON(code, H{"error": err})
}
