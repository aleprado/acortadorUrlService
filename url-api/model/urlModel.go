package model

import "time"

type ShortenedURL struct {
	Hash      string    `json:"hash"`
	Original  string    `json:"original"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"` 
}
