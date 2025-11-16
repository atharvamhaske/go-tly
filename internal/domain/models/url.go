package models

import (
	"time"
)

type URL struct {
	ID            string             `bson:"_id,omitempty" json:"id"`
	ShortKey      string             `bson:"short_key" json:"short_key"`
	LongURL       string             `bson:"long_url" json:"long_url"`
	UserID        *string            `bson:"user_id,omitempty" json:"user_id"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
    Expiry        *time.Time         `bson:"expiry,omitempty" json:"expiry"`
    DomainHealth  *DomainHealth      `bson:"domain_health,omitempty" json:"domain_health"`
    Clicks        int64              `bson:"clicks" json:"clicks"`
    IsActive      bool               `bson:"is_active" json:"is_active"`
}

// will add custom alias in next version

type ShortenRequest struct {
	LongURL string `json:"long_url,omitempty"`
	Expiry  *string `json:"expiry,omitempty"`
}

type ShortenResponse struct {
    ShortURL string `json:"short_url"`
    ShortKey string `json:"short_key"`
}
