package main

import (
	"net/http"

	"geek/web"
)

// * Version Without Context
/*
```go
func appWithoutContext() {
	app := web.NewEngine()
	app.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})
	app.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%s] = %q\n", k, v)
		}
	})
	app.Run(":9001")
}
```
*/

func appWithContext() {
	app := web.NewEngine()
	app.GET("/", func(ctx *web.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello Geek!</h1>")
	})
	app.GET("/hello", func(ctx *web.Context) {
		ctx.String(http.StatusOK, "Hello, %s\n", ctx.Query("name"))
	})
	app.POST("/login", func(ctx *web.Context) {
		ctx.JSON(http.StatusOK, web.H{
			"username": ctx.PostForm("username"),
			"password": ctx.PostForm("password"),
		})
	})
	app.Run(":9001")
}

func main() {
	appWithContext()
}
