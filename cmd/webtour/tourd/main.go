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

	app.GET("/hello/:name", func(ctx *web.Context) {
		ctx.JSON(http.StatusOK, web.H{
			"name": ctx.Param("name"),
			"path": ctx.Path,
		})
	})

	app.Run(":9001")
}

func appWithGroup() {
	app := web.NewEngine()

	// Enable Request Logger
	app.Use(web.Logger())

	app.GET("/", func(ctx *web.Context) {
		ctx.HTML(http.StatusOK, "<h1>Home Index</h1>")
	})
	app.GET("/index", func(ctx *web.Context) {
		ctx.HTML(http.StatusOK, "<h1>Home Index</h1>")
	})
	app.GET("/apis", func(ctx *web.Context) {
		ctx.JSON(http.StatusOK, app.Router())
	})

	// * V1 Router Group
	v1 := app.Group("/v1")
	v1.GET("/", func(ctx *web.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello Geek!</h1>")
	})
	v1.GET("/hello", func(ctx *web.Context) {
		ctx.String(http.StatusOK, "Hello, %s\n", ctx.Query("name"))
	})

	// * V2 Router Group
	v2 := app.Group("/v2")
	v2.POST("/login", func(ctx *web.Context) {
		ctx.JSON(http.StatusOK, web.H{
			"username": ctx.PostForm("username"),
			"password": ctx.PostForm("password"),
		})
	})
	v2.GET("/hello/:name", func(ctx *web.Context) {
		ctx.JSON(http.StatusOK, web.H{
			"name": ctx.Param("name"),
			"path": ctx.Path,
		})
	})

	// * Nested Router Group
	static := v2.Group("/static")
	static.GET("/", func(ctx *web.Context) {
		ctx.HTML(http.StatusOK, "<h2>Home</h2")
	})
	static.GET("/home", func(ctx *web.Context) {
		ctx.HTML(http.StatusOK, "<h2>Home</h2")
	})

	app.Run(":9001")
}

func main() {
	appWithGroup()
}
