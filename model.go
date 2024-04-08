package main

import "time"

// APIKey reflects a hash in Redis
type APIKey struct {
	ID        int       `json:"id"`
	Key       string    `json:"key"`
	UserID    string    `json:"user_id,omitempty"`
	CreatedAT time.Time `json:"created_at"`
}

// User reflects a hash in Redis
type User struct {
	ID              int       `json:"id"`
	Username        string    `json:"username"`
	CurrentAPIKeyID int       `json:"current_api_key_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// URL reflects a ZSet (sorted set) in Redis
type URL struct {
	OriginalURL string    `json:"original_url"`
	Shortcode   string    `json:"shortcode"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}
