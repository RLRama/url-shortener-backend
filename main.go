package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"os"
)

var ctx = context.Background()
var pepper = os.Getenv("PEPPER")

func init() {
	err := loadEnv()
	if err != nil {
		return
	}
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
	app := iris.Default()

	// test handlers (will be dropped soon)
	app.Post("/hello-world", helloWorldTest)
	app.Post("/redis-test", redisTest)

	// user handlers
	user := app.Party("/user")
	{
		user.Post("/register", registerUserHandler)
	}

	return app
}
