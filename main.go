package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"sync"
)

type WrapperApplication struct {
	*iris.Application
}

var applicationPool = sync.Pool{New: func() interface{} {
	return &WrapperApplication{}
}}

func newWrapperApplication() *WrapperApplication {
	app := applicationPool.Get().(*WrapperApplication)
	// set the application to the original one in order to have access to iris's implementation.
	app.Application = iris.New()
	return app
}

func (app *WrapperApplication) CustomHandle(method string, relativePath string, fns ...context.Handler) {
	handlers := make([]iris.Handler, 0)
	for i, _ := range fns {
		handlers = append(handlers, wrapper(fns[i]))
	}
	app.Handle(method, relativePath, handlers...)
}

// wrapper handler with generator context next
func wrapper(handle context.Handler) iris.Handler {
	return func(ctx iris.Context) {
		handle(ctx)
		ctx.Next()
	}
}

func main() {
	app := newWrapperApplication()

	// Register a view engine on .html files inside the ./view/** directory.
	app.RegisterView(iris.HTML("./view", ".html"))

	// this will be executed with generator context next
	app.CustomHandle("GET", "/hi/{firstname:alphabetical}", func(ctx iris.Context) {
		firstname := ctx.Params().GetString("firstname")

		ctx.ViewData("firstname", firstname)
		ctx.Gzip(true)

		ctx.View("hi.html")
	})

	app.DoneGlobal(after)

	app.Run(iris.Addr(":8080"))
}

func after(ctx iris.Context) {
	fmt.Println("executed done global handler !")
}
