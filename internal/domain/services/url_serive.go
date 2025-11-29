package services

import (
	"context"
	"errors"
	"time"

	"github.com/atharvamhaske/go-tly/internal/domain"
	"github.com/atharvamhaske/go-tly/internal/domain/models"
	"github.com/google/uuid"
)

type URLService struct {
	urlRepo domain.URLRepo
	cache   domain.URLCache
}

func NewURLService(urlRepo domain.URLRepo, cache domain.URLCache) domain.URLService {
	return &URLService{
		urlRepo: urlRepo,
		cache:   cache,
	}
}

// Shorten creates a new shortened URL and returns its key and short URL representation.
func (s *URLService) Shorten(ctx context.Context, req models.ShortenRequest) (*models.ShortenResponse, error) {
	if req.LongURL == "" {
		return nil, errors.New("long_url is required")
	}

	now := time.Now()

	var expiry *time.Time
	if req.Expiry != nil {
		// Interpret expiry as duration string like "24h".
		d, err := time.ParseDuration(*req.Expiry)
		if err != nil {
			return nil, err
		}
		e := now.Add(d)
		expiry = &e
	}

	shortKey := uuid.NewString()[:8]

	url := &models.URL{
		ShortKey:  shortKey,
		LongURL:   req.LongURL,
		CreatedAt: now,
		Expiry:    expiry,
		IsActive:  true,
	}

	if err := s.urlRepo.Save(ctx, url); err != nil {
		return nil, err
	}

	// Best-effort cache set; ignore errors.
	_ = s.cache.SetKey(ctx, shortKey, req.LongURL)

	// The actual domain/base URL can be added at the handler layer; here we just return the key.
	resp := &models.ShortenResponse{
		ShortURL: shortKey,
		ShortKey: shortKey,
	}

	return resp, nil
}

// Edit updates mutable fields of an existing short URL.
func (s *URLService) Edit(ctx context.Context, key string, req models.EditURLRequest) error {
	if key == "" {
		return errors.New("key is required")
	}

	updates := make(map[string]any)
	if req.LongURL != "" {
		updates["long_url"] = req.LongURL
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	if err := s.urlRepo.UpdateByKey(ctx, key, updates); err != nil {
		return err
	}

	// Keep cache in sync; ignore cache errors.
	if req.LongURL != "" {
		_ = s.cache.SetKey(ctx, key, req.LongURL)
	}

	return nil
}

// Delete removes a short URL and evicts it from cache.
func (s *URLService) Delete(ctx context.Context, key string) error {
	if key == "" {
		return errors.New("key is required")
	}

	if err := s.urlRepo.DeleteByKey(ctx, key); err != nil {
		return err
	}

	// Best-effort cache delete.
	_, _ = s.cache.DeleteKey(ctx, key)

	return nil
}

// Resolve returns the original long URL for a given short key, using cache when possible.
func (s *URLService) Resolve(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", errors.New("key is required")
	}

	// First try cache.
	if longURL, err := s.cache.GetKey(ctx, key); err == nil && longURL != "" {
		return longURL, nil
	}

	// Fallback to repository.
	url, err := s.urlRepo.FindByKey(ctx, key)
	if err != nil {
		return "", err
	}
	if url == nil || !url.IsActive {
		return "", errors.New("url not found or inactive")
	}

	// Check expiry, if set.
	if url.Expiry != nil && time.Now().After(*url.Expiry) {
		return "", errors.New("url expired")
	}

	// Refresh cache; ignore errors.
	_ = s.cache.SetKey(ctx, key, url.LongURL)

	return url.LongURL, nil
}
