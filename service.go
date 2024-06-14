package main

import (
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"unicode"
)

func validateFieldLength(field, fieldName string, min, max int) error {
	length := len(field)
	if length < min || length > max {
		return fmt.Errorf("%s must be between %d and %d characters", fieldName, min, max)
	}
	return nil
}

func sanitizeField(field string) string {
	// replace characters not inside the following regex
	reg := regexp.MustCompile("[^a-zA-Z0-9_]+")
	return reg.ReplaceAllString(field, "")
}

func hasAllCharTypes(field string) error {
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	for _, char := range field {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("field contains no uppercase characters")
	}
	if !hasLower {
		return errors.New("field contains no lowercase characters")
	}
	if !hasNumber {
		return errors.New("field contains no number")
	}
	if !hasSpecial {
		return errors.New("field contains no special characters")
	}

	return nil
}

func checkUsernameExists(username string) (bool, error) {
	nextUserId, err := rdb.Get(ctx, "next_user_id").Uint64()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, fmt.Errorf("error getting next user id from redis: %w", err)
	}

	for i := uint64(0); i < nextUserId; i++ {
		userKey := fmt.Sprintf("user:%d", i)
		storedUsername, err2 := rdb.HGet(ctx, userKey, "username").Result()
		if err2 != nil && !errors.Is(err2, redis.Nil) {
			return false, fmt.Errorf("error checking username for user %d: %w", i, err)
		}
		if storedUsername == username {
			return true, nil
		}
	}

	return false, nil
}

func hashPassword(password string) (string, error) {
	pepperedPass := []byte(password + pepper)

	hash, err := bcrypt.GenerateFromPassword(pepperedPass, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}

	return string(hash), nil
}

func verifyPassword(password, hash string) bool {
	pepperedPass := []byte(password + pepper)
	err := bcrypt.CompareHashAndPassword([]byte(hash), pepperedPass)
	return err == nil
}
