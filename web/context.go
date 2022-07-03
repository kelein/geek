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
}

// newContext contains http request and response info
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Request: req,
		Writer:  w,
		Method:  req.Method,
		Path:    req.URL.Path,
	}
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
func (ctx *Context) HTML(code int, html string) {
	ctx.Status(code)
	ctx.SetHeader("Content-Type", "text/html")
	ctx.Writer.Write([]byte(html))
}
