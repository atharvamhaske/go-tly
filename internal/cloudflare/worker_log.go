package cloudflare

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/atharvamhaske/go-tly/internal/domain/models"
)

// WorkerLogVerifier verifies Cloudflare Worker log signatures.
// This is a utility helper used by Analytics ingester to verify
// that click events from Cloudflare Workers are authentic.
type WorkerLogVerifier struct {
	secret string
}

// NewWorkerLogVerifier creates a new log verifier.
func NewWorkerLogVerifier(secret string) *WorkerLogVerifier {
	return &WorkerLogVerifier{secret: secret}
}

// VerifyLogEntry verifies a log entry from Cloudflare Worker.
// Returns the parsed ClickEvents if signature is valid.
func (v *WorkerLogVerifier) VerifyLogEntry(timestamp, signature, body string) (*models.ClickEvents, error) {
	// Check timestamp freshness (within 5 minutes)
	ts, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp: %w", err)
	}

	if time.Since(ts) > 5*time.Minute {
		return nil, fmt.Errorf("timestamp too old")
	}

	// Verify signature
	if !v.VerifySignature(timestamp, signature, body) {
		return nil, fmt.Errorf("invalid signature")
	}

	// Parse body into ClickEvents
	var evt models.ClickEvents
	if err := json.Unmarshal([]byte(body), &evt); err != nil {
		return nil, fmt.Errorf("parse body: %w", err)
	}

	return &evt, nil
}

// VerifySignature verifies HMAC signature.
func (v *WorkerLogVerifier) VerifySignature(timestamp, signature, body string) bool {
	mac := hmac.New(sha256.New, []byte(v.secret))
	mac.Write([]byte(timestamp + body))
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}

