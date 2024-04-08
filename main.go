package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

func main() {

	err3 := godotenv.Load()
	if err3 != nil {
		log.Fatal("Error loading .env file")
	}

	redisURI := os.Getenv("REDIS_URI")

	addr, err2 := redis.ParseURL(redisURI)
	if err2 != nil {
		log.Fatal(err2)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err3 = r.Run()
	if err3 != nil {
		return
	}
}
