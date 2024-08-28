package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"os"
)

var ctx = context.Background()
var pepper = os.Getenv("PEPPER")
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

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
	app.Get("/protected", authMiddleware, func(ctx iris.Context) {
		_, err := ctx.WriteString("You accessed the protected route")
		if err != nil {
			return
		}
	})

	// user handlers
	user := app.Party("/user")
	{
		user.Post("/register", handleUserRegistration)
		user.Post("/login", handleLogin)
	}
	authenticatedUser := app.Party("/user")
	authenticatedUser.Use(authMiddleware)
	{
		// update username, password, drop user, etc
		authenticatedUser.Put("/update-password", handleUpdatePassword)
		authenticatedUser.Put("/update-username", handleUpdateUsername)
		authenticatedUser.Post("/logout", handleLogout)
	}

	return app
}
