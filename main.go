package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

var rdb *redis.Client
var ctx = context.Background()

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisURI := os.Getenv("REDIS_URI")

	opt, err2 := redis.ParseURL(redisURI)
	if err2 != nil {
		log.Fatal(err2)
	}

	rdb = redis.NewClient(opt)

}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(ErrorHandlingMiddleware())
	r.Use(RateLimitMiddleware(1, 1))
	r.Use(RequestLoggingMiddleware())

	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", RegisterUserHandler)
	}

	return r

}

func main() {

	r := setupRouter()

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
