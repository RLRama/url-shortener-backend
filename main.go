package main

import (
	"context"
	"github.com/kataras/iris/v12"
)

var ctx = context.Background()

func init() {
	connectToDatabase()
}

func main() {
	app := newApp()

	err := app.Listen(":8080")
	if err != nil {
		return
	}
}

func newApp() *iris.Application {
	app := iris.New()

	// here go the routes
	app.Post("/hello-world", helloWorldTest)

	return app
}
