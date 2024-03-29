package main

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"time"

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
		ctx.HTML(http.StatusOK, "<h1>Hello Geek!</h1>", nil)
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

	// Enable Recover
	app.Use(web.Recover())

	// * Eanble HTML Template Render
	dir, _ := os.Getwd()
	app.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormateAsDate,
	})
	app.LoadGlobHTML(path.Join(dir, "assets", "templates/*"))

	// * Eanble Static File Support
	app.Static("/assets", path.Join(dir, "assets"))

	app.GET("/", func(ctx *web.Context) {
		ctx.HTML(http.StatusOK, "<h1>Home Index</h1>", nil)
	})
	app.GET("/index", func(ctx *web.Context) {
		ctx.HTML(http.StatusOK, "<h1>Home Index</h1>", nil)
	})
	app.GET("/apis", func(ctx *web.Context) {
		ctx.JSON(http.StatusOK, app.Router())
	})

	// * V1 Router Group
	v1 := app.Group("/v1")
	v1.GET("/", func(ctx *web.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello Geek!</h1>", nil)
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
		ctx.HTML(http.StatusOK, "<h2>Home</h2", nil)
	})
	static.GET("/home", func(ctx *web.Context) {
		ctx.HTML(http.StatusOK, "<h2>Home</h2", nil)
	})
	static.GET("/except", func(ctx *web.Context) {
		ctx.JSON(http.StatusGatewayTimeout, web.H{
			"error": http.ErrHandlerTimeout.Error(),
			"code":  http.StatusGatewayTimeout,
		})
	})

	// * Test Template Render
	static.GET("/demo", func(ctx *web.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	// * Test Recover
	static.GET("/panic", func(ctx *web.Context) {
		v := []string{"test", "panic"}
		ctx.String(http.StatusOK, v[3])
	})

	app.Run(":9001")
}

// FormateAsDate for html render
func FormateAsDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func main() {
	appWithGroup()
}
