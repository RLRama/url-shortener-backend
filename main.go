package main

import "github.com/kataras/iris/v12"

func main() {
}

func newApp() *iris.Application {
	app := iris.New()

	app.Logger().SetLevel("debug")
	app.Use(setAllowedResponses)
}

func setAllowedResponses(ctx iris.Context) {

	ctx.Negotiation().JSON().XML().YAML().MsgPack()
	ctx.Negotiation().Accept.JSON()
	ctx.Next()

}
