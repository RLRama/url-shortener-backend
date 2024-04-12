package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

var rdb *redis.Client
var ctx = context.Background()

func ShortenURL(c *gin.Context) {
	var requestBody URL
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortenedKey := generateUniqueKey()

	err := rdb.HSet(ctx, fmt.Sprint("url:%s", shortenedKey), map[string]interface{}{
		"original_url":  requestBody.OriginalURL,
		"creation_time": time.Now().Unix(),
	}).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"shortened_url": shortenedKey})
}

func GenerateAPIKey(c *gin.Context) {
	var requestBody struct {
		OwnerID string `json:"owner_id"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	apiKey := generateUniqueKey()

	err := rdb.HSet(ctx, fmt.Sprint("apikey:%s", apiKey), map[string]interface{}{
		"owner_id": requestBody.OwnerID,
		"api"
	})
}
