package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	redisURI := "rediss://default:AVNS_3T4ZlPmnGar1bWRDsXz@redis-url-shortener-rlri.a.aivencloud.com:14949"

	addr, err := redis.ParseURL(redisURI)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(addr)

	err = rdb.Set(ctx, "key", "hello world", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("The value of key is:", val)
}
