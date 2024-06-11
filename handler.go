package main

import (
	"github.com/kataras/iris/v12"
)

func helloWorldTest(ctx iris.Context) {
	_, err := ctx.WriteString("Hello World")
	if err != nil {
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
}

func redisTest(ctx iris.Context) {
	err := rdb.Set(ctx, "key0", "dickson", 0).Err()
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	value, err4 := rdb.Get(ctx, "key0").Result()
	if err4 != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err4)
		return
	}

	err6 := ctx.JSON(value)
	if err6 != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err6)
		return
	}

}
