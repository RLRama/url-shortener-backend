package main

type APIKey struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	CreationTime int64  `json:"creation_time"`
}

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	CreationTime int64  `json:"creation_time"`
}

type URL struct {
	ID           string `json:"id"`
	OriginalURL  string `json:"original_url"`
	CreatorID    string `json:"creator_id"`
	CreationTime int64  `json:"creation_time"`
}
