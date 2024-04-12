package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
)

var rdb *redis.Client
var ctx = context.Background()

func init() {

	err0 := godotenv.Load()
	if err0 != nil {
		log.Fatal("Error loading .env file")
	}

	redisURI := os.Getenv("REDIS_URI")

	opt, err2 := redis.ParseURL(redisURI)
	if err2 != nil {
		log.Fatal(err2)
	}

	rdb = redis.NewClient(opt)

}

func main() {
	r := gin.Default()

	r.Use(AuthMiddleware())

	r.GET("/getUrls", GetUserShortenedURLs)
	r.POST("/postUrl", CreateShortenedURL)
	r.POST("/generateApiKey", func(c *gin.Context) {
		userID := c.PostForm("user_id")

		apiKey, err := StoreAPIKey(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to generate API Key"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"api_key": apiKey})
	})

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
