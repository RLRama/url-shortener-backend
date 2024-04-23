package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			fmt.Println(c.Errors.String())
			c.JSON(http.StatusInternalServerError, gin.H{"error": c.Errors.String()})
			c.Abort()
			return
		}
	}
}

func RateLimitMiddleware(rateLimit int, burstLimit int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(rateLimit), burstLimit)

	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
		} else {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			c.Abort()
		}
	}
}

func RequestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()

		fmt.Printf("[%s] %s %s %d %v %s\n", end.Format("2006/01/02 - 15:04:05"), clientIP, method, statusCode, latency, path)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header"})
			c.Abort()
			return
		}

		user, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}

func ValidateToken(tokenString string) (*User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte("secret"), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, errors.New("username not found in token claims")
	}

	user, err := FindUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func FindUserByUsername(username string) (*User, error) {
	keys, err := rdb.Keys(ctx, "user:*").Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		userData, err := rdb.HGetAll(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		if userData["username"] == username {
			user := &User{
				Username: userData["username"],
				Password: userData["password"],
			}
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}
