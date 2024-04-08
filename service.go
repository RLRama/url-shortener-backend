package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

var ctx = context.Background()

func GenerateAPIKey(redisClient *redis.Client) (*APIKey, error) {
	apiKeyBytes, err := rand.Read(20)
	if err != nil {
		return nil, errors.New("failed to generate random bytes for api key")
	}

	apiKey := base64.StdEncoding.EncodeToString(apiKeyBytes)

	hashedKey, err := bcrypt.GenerateFromPassword(apiKey, bcrypt.DefaultCost)
}



func ValidateAPIKey(key string) (*APIKey, error) {}()  {

}