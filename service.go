package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

// URL functions

func createURLHandler(c *gin.Context) {
	var urlData URL
	if err := c.ShouldBindJSON(&urlData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user id not found"})
		return
	}

	urlData.CreatorID = userID.(string)

	if urlData.OriginalURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing URL or Creator ID"})
		return
	}

	urlShortcode, err := generateUniqueURLShortcode()
	if err != nil {
		rdb.Del(ctx, urlShortcode)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if urlShortcode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate unique URL short code"})
	}

	urlData.CreationTime = time.Now().Unix()

	key := "url:" + urlShortcode
	fields := make(map[string]interface{})
	fields["original_url"] = urlData.OriginalURL
	fields["creator_id"] = userID
	fields["creation_time"] = urlData.CreationTime

	if err := rdb.HSet(ctx, key, fields).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"url_id": urlShortcode})
}

func generateUniqueURLShortcode() (string, error) {
	var urlShortcode string
	for {
		urlShortcode = generateRandomString(10)

		key := "url:" + urlShortcode
		exists, err := rdb.Exists(ctx, key).Result()
		if err != nil {
			return "", err
		}

		if exists == 0 {
			break
		}
	}
	return urlShortcode, nil
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

// Auth functions

func authMiddleware(c *gin.Context) {
	apiKey := c.GetHeader("x-api-key")
	if apiKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing API key"})
		c.Abort()
		return
	}

	keyExists, err := rdb.Exists(ctx, "apiKey:"+apiKey).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if keyExists == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		c.Abort()
		return
	}

	userID, err := getUserIDByAPIKey(apiKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.Set("userID", userID)

	c.Next()
}

func getUserIDByAPIKey(apiKey string) (string, error) {
	userID, err := rdb.HGet(ctx, "apikey:"+apiKey, "user_id").Result()
	if err != nil {
		return "", err
	}

	if userID == "" {
		return "", errors.New("invalid API key")
	}

	return userID, nil
}

// User functions

func registerUser(username string) (string, error) {
	userID, err := generateUserID()
	if err != nil {
		return "", err
	}

	newUser := User{
		ID:           userID,
		Username:     username,
		CreationTime: time.Now().Unix(),
	}

	userKey := "user:" + userID

	if err := rdb.HSet(ctx, userKey, map[string]interface{}{
		"username":      newUser.Username,
		"creation_time": newUser.CreationTime,
	}).Err(); err != nil {
		return "", err
	}

	apiKey, err := generateAPIKey()
	if err != nil {
		rdb.Del(ctx, userKey)
		return "", err
	}

	apiKeyKey := "apiKey:" + apiKey

	if err := rdb.HSet(ctx, apiKeyKey, map[string]interface{}{
		"user_id":       userID,
		"creation_time": time.Now().Unix(),
	}).Err(); err != nil {
		rdb.Del(ctx, userKey)
		return "", err
	}

	return apiKey, nil
}

func generateUserID() (string, error) {
	var userID string
	for {
		userID = generateRandomString(20)

		key := "user" + userID
		exists, err := rdb.Exists(ctx, key).Result()
		if err != nil {
			return "", err
		}
		if exists == 0 {
			break
		}
	}

	return userID, nil
}

// Utilities

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}
	for i := range randomBytes {
		randomBytes[i] = charset[int(randomBytes[i])%len(charset)]
	}
	return string(randomBytes)
}

// API Key functions

func generateAPIKey() (string, error) {

	const apiKeyLength = 32
	var apiKey string

	for {

		apiKey = generateRandomString(apiKeyLength)

		key := "apikey:" + apiKey
		exists, err := rdb.Exists(ctx, key).Result()
		if err != nil {
			return "", err
		}
		if exists == 0 {
			break
		}
	}

	return apiKey, nil
}

func getOrCreateAPIKeyForUser(userID, providedAPIKey string) (string, error) {
	existingUserID, err := rdb.HGet(ctx, "apikey:"+providedAPIKey, "user_id").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}

}
