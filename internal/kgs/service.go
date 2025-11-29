package kgs

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// Service defines the business logic for generating keys and interacting with
// backing stores like Redis/Mongo (hooked in via repositories).
type Service interface {
	GenerateKey(ctx context.Context) (string, error)
}

// simpleService is a minimal implementation that just generates random keys.
// You can extend this later to use Redis queues and Mongo backup.
type simpleService struct{}

func NewService() Service {
	return &simpleService{}
}

func (s *simpleService) GenerateKey(ctx context.Context) (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("generate key: %w", err)
	}
	// URL-safe base64, strip padding and maybe trim for length.
	key := base64.RawURLEncoding.EncodeToString(b[:])
	return key, nil
}


