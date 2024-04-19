package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
