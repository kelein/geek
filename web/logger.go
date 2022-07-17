package web

import (
	"fmt"
	"log"
	"time"
)

// Terminal Colors
const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Color represents a text color.
type Color uint8

// Add adds the color to any type
func (c Color) Add(s interface{}) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", uint8(c), s)
}

// Bold adds a bold color to the given value
func (c Color) Bold(s interface{}) string {
	return fmt.Sprintf("\x1b[1;%dm%v\x1b[0m", uint8(c), s)
}

func statusColor(code int) string {
	if code >= 200 && code < 400 {
		return Green.Bold(code)
	}
	if code >= 400 && code < 500 {
		return Red.Bold(code)
	}
	return Magenta.Bold(code)
}

// Logger for default request log middleware
func Logger() HandlerFunc {
	return func(ctx *Context) {
		start := time.Now()
		ctx.Next()
		elaps := time.Since(start)
		log.Printf("%s - %s - %s - %s - %s", ctx.Request.RemoteAddr, statusColor(ctx.StatusCode), ctx.Request.Method, ctx.Request.RequestURI, elaps)
	}
}
