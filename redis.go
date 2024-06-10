package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

var rdb *redis.Client

func connectToDatabase() *redis.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectionString := os.Getenv("REDIS_URI")

	opt, err2 := redis.ParseURL(connectionString)
	if err2 != nil {
		log.Fatal(err2)
	}

	rdb = redis.NewClient(opt)

	_, err3 := rdb.Ping(context.Background()).Result()
	if err3 != nil {
		log.Fatal(err3)
	}

	return rdb
}
