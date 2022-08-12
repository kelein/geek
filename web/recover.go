package web

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// Recover -
func Recover() HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf("%s", err)
				log.Printf("%s", Trace(msg))
				ctx.Fail(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}
		}()

		ctx.Next()
	}
}

// Trace print stack for debug
func Trace(msg string) string {
	pcs := [32]uintptr{}
	n := runtime.Callers(3, pcs[:])

	s := strings.Builder{}
	s.WriteString(msg + "\nTraceback:")

	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		s.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return s.String()
}
