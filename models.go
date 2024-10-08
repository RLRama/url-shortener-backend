package main

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// User struct
type User struct {
	ID        uint64    `json:"id" redis:"id"`
	Username  string    `json:"username" redis:"username"`
	Password  string    `json:"-" redis:"password"`
	CreatedAt time.Time `json:"created_at" redis:"created_at"`
	UpdatedAt time.Time `json:"updated_at" redis:"updated_at"`
}

// URL struct
type URL struct {
	ID          uint64    `json:"id" redis:"id"`
	OriginalURL string    `json:"original_url" redis:"original_url"`
	ShortURL    string    `json:"short_url" redis:"short_url"`
	CreatedAt   time.Time `json:"created_at" redis:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" redis:"updated_at"`
	UserID      uint64    `json:"user_id" redis:"user_id"`
	ViewCount   uint64    `json:"view_count" redis:"view_count"`
}

// RegisterUserRequest for registering a user
type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// PasswordUpdateRequest for changing a user's password
type PasswordUpdateRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

// UsernameUpdateRequest for changing a user's username
type UsernameUpdateRequest struct {
	NewUsername string `json:"new_username" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

// Claims Login claims struct
type Claims struct {
	Username string `json:"username"`
	jwt.MapClaims
}
