package main

// APIKey reflects the API key model as a hash in Redis
type APIKey struct {
	OwnerID      string `json:"owner_id"`
	APIKey       string `json:"api_key"`
	CreationTime int64  `json:"creation_time"`
}

// URL reflects the URL model as a hash in Redis
type URL struct {
	OriginalURL  string `json:"original_url"`
	CreationTime int64  `json:"creation_time"`
	CreatorID    string `json:"creator_id"`
}

type User struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	CreationTime int64  `json:"creation_time"`
}
