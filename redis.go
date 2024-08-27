package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

var rdb *redis.Client

func loadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	return nil
}

func connectToDatabase() *redis.Client {

	connectionString := os.Getenv("REDIS_URI")

	opt, err2 := redis.ParseURL(connectionString)
	if err2 != nil {
		log.Fatal(err2)
	}

	rdb = redis.NewClient(opt)

	_, err3 := rdb.Ping(context.Background()).Result()
	if err3 != nil {
		panic(err3)
	}

	_, err4 := rdb.Ping(context.Background()).Result()
	if err4 != nil {
		log.Fatal(err4)
	}

	return rdb
}
