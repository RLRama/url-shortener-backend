package main

import (
	"github.com/kataras/iris/v12"
)

func newApp() *iris.Application {
	app := iris.New()

	app.Post("/", func(ctx iris.Context) {})
}

func main() {
}
