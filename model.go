package main

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Salt     string `json:"salt"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

type ShortenedURL struct {
	Original  string `json:"original"`
	Shortcode string `json:"shortcode"`
	UserID    string `json:"userId"`
	Clicks    int    `json:"clicks"`
	CreatedAt int64  `json:"createdAt"`
}
