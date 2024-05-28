package main

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Salt      string `json:"salt"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type URL struct {
	ID          int    `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
	CreatedAt   string `json:"created_at"`
	ExpiresAt   string `json:"expires_at"`
	UserID      int    `json:"user_id"`
}
