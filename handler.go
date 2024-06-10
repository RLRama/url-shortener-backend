package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
)

func helloWorldTest(ctx iris.Context) {
	_, err := ctx.WriteString("Hello World")
	if err != nil {
		return
	}
}

func redisTest(ctx iris.Context) {
	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		_, err2 := ctx.WriteString(err.Error())
		if err2 != nil {
			return
		}
		return
	}

	value, err4 := rdb.Get(ctx, "key").Result()
	if err4 != nil {
		_, err3 := ctx.WriteString(err4.Error())
		if err3 != nil {
			return
		}
		return
	}

	_, err5 := ctx.WriteString(fmt.Sprintf("value: %s", value))
	if err5 != nil {
		return
	}
}
