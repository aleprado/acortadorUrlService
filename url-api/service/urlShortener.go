package service

import (
	"acortadorUrlService/components/database"
	"acortadorUrlService/url-api/model"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"time"
)

type UrlShortener struct {
	DB *database.DDBClient
}

func NewUrlShortener(db *database.DDBClient) *UrlShortener {
	return &UrlShortener{DB: db}
}

func (s *UrlShortener) ShortenOrFetch(ctx context.Context, originalURL string) (*model.ShortenedURL, error) {
	hash := generateHash(originalURL)

	existingOriginal, err := s.DB.GetURL(ctx, hash)
	if err == nil && existingOriginal != "" {
		return &model.ShortenedURL{
			Hash:     hash,
			Original: existingOriginal,
		}, nil
	}
	now := time.Now()
	newShort := &model.ShortenedURL{
		Hash:      hash,
		Original:  originalURL,
		CreatedAt: now,
		ExpiresAt: now.Add(7 * 24 * time.Hour),
	}

	err = s.DB.SaveURL(ctx, hash, originalURL)
	if err != nil {
		return nil, err
	}

	return newShort, nil
}

func (s *UrlShortener) Delete(ctx context.Context, originalURL string) error {
	hash := generateHash(originalURL)
	return s.DB.DeleteURL(ctx, hash)
}

//TODO: move this function to utils 
func generateHash(url string) string {
	h := sha256.Sum256([]byte(url))
	encoded := base64.URLEncoding.EncodeToString(h[:])
	return strings.TrimRight(encoded[:8], "=")
}

func (s *UrlShortener) GetOriginalUrl(ctx context.Context, hash string) (string, error) {
	return s.DB.GetURL(ctx, hash)
}
