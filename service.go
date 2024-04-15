package main

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func createURLHandler(c *gin.Context) {
	var urlData URL
	if err := c.ShouldBindJSON(&urlData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if urlData.OriginalURL == "" || urlData.CreatorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing URL or Creator ID"})
		return
	}

	urlShortcode := generateURLShortcode()

	urlData.CreationTime = time.Now().Unix()

	key := "url:" + urlShortcode
	fields := make(map[string]interface{})
	fields["original_url"] = urlData.OriginalURL
	fields["creator_id"] = urlData.CreatorID
	fields["creation_time"] = urlData.CreationTime

	if err := rdb.HSet(ctx, key, fields, urlData.OriginalURL).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"url_id": urlShortcode})
}

func generateURLShortcode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		panic(err)
	}

	shortcode := base64.URLEncoding.EncodeToString(randomBytes)

	if len(shortcode) < 10 {
		padding := make([]byte, 10-len(shortcode))
		shortcode += string(padding)
	} else if len(shortcode) > 10 {
		shortcode = shortcode[:10]
	}

	for i := range shortcode {
		if shortcode[i] == '=' {
			shortcode = shortcode[:i] + charset[i%len(charset):i%len(charset)+1] + shortcode[i+1:]
		}
	}

	return shortcode
}
