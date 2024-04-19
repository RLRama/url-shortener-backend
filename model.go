package main

import (
	"golang.org/x/time/rate"
	"sync"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ShortenedURL struct {
	Original  string `json:"original"`
	Shortcode string `json:"shortcode"`
	UserID    string `json:"userId"`
	Clicks    int    `json:"clicks"`
	CreatedAt int64  `json:"createdAt"`
}

type RateLimiter struct {
	limiter *rate.Limiter
	mu      sync.Mutex
}
