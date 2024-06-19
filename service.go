package main

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"time"
	"unicode"
)

// works for username as well as password
func validateFieldLength(field, fieldName string, min, max int) error {
	length := len(field)
	if length < min || length > max {
		return fmt.Errorf("%s must be between %d and %d characters", fieldName, min, max)
	}
	return nil
}

func validatePasswordCharTypes(field string) error {
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

func getNextUserId() (uint64, error) {
	nextUserId, err := rdb.Get(ctx, "next_user_id").Uint64()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, fmt.Errorf("error getting next user id from redis: %w", err)
	}

	return nextUserId, nil
}

func checkUsernameExists(username string) (bool, error) {
	nextUserId, err := getNextUserId()
	if err != nil {
		return false, fmt.Errorf("error getting next user id from redis: %w", err)
	}

	for i := uint64(0); i < nextUserId; i++ {
		userKey := fmt.Sprintf("user:%d", i)
		storedUsername, err2 := rdb.HGet(ctx, userKey, "username").Result()
		if err2 != nil && !errors.Is(err2, redis.Nil) {
			return false, fmt.Errorf("error checking username for user %d: %w", i, err2)
		}
		if storedUsername == username {
			return true, nil
		}
	}

	return false, nil
}

func getUserByUsername(username string) (*User, error) {
	nextUserId, err := getNextUserId()
	if err != nil {
		return nil, fmt.Errorf("error getting next user id from redis: %w", err)
	}

	for i := uint64(0); i <= nextUserId; i++ {
		userKey := fmt.Sprintf("user:%d", i)
		storedUsername, err2 := rdb.HGet(ctx, userKey, "username").Result()
		if err2 != nil && !errors.Is(err2, redis.Nil) {
			return nil, fmt.Errorf("error checking username for user %d: %w", i, err2)
		}
		if storedUsername == username {
			userData, err3 := rdb.HGetAll(ctx, userKey).Result()
			if err3 != nil && !errors.Is(err3, redis.Nil) {
				return nil, fmt.Errorf("error getting user data for user %d: %w", i, err3)
			}

			createdAt, err4 := time.Parse(time.RFC3339, userData["created_at"])
			if err4 != nil {
				return nil, fmt.Errorf("error parsing created_at for user %d: %w", i, err4)
			}

			updatedAt, err5 := time.Parse(time.RFC3339, userData["updated_at"])
			if err5 != nil {
				return nil, fmt.Errorf("error parsing updated_at for user %d: %w", i, err5)
			}

			user := &User{
				ID:        i,
				Username:  userData["username"],
				Password:  userData["password"],
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			}
			return user, nil
		}
	}

	return nil, errors.New("user not found")
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

func generateToken(username string) (string, error) {
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := Claims{
		Username: username,
		MapClaims: jwt.MapClaims{
			"exp": expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
