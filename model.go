package main

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
