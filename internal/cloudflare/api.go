package cloudflare

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/atharvamhaske/go-tly/internal/domain/models"
)

type HealthChecker struct {
	httpClient *http.Client
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: false,
				},
			},
		},
	}
}

// CheckDomainHealth validates a URI's health.
func (h *HealthChecker) CheckDomainHealth(ctx context.Context, uri string) (*models.DomainHealth, error) {
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	health := &models.DomainHealth{
		DNSResolved: false,
		SSLValid:    false,
		StatusCode:  0,
		LatencyMS:   0,
	}

	// Check DNS resolution
	start := time.Now()
	resp, err := h.httpClient.Get(uri)
	latency := time.Since(start)
	health.LatencyMS = latency.Milliseconds()

	if err != nil {
		health.Message = fmt.Sprintf("DNS/Connection error: %v", err)
		return health, nil // Return partial health info
	}
	defer resp.Body.Close()

	health.DNSResolved = true
	health.StatusCode = resp.StatusCode

	// Check SSL if HTTPS
	if parsedURL.Scheme == "https" {
		if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
			health.SSLValid = true
		}
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		health.Message = "Domain is healthy"
	} else {
		health.Message = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}

	return health, nil
}