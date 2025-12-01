package cloudflare

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/atharvamhaske/go-tly/internal/domain"
	"github.com/atharvamhaske/go-tly/internal/domain/models"
)

// cloudflareHealthService implements domain.HealthService interface.
type cloudflareHealthService struct {
	httpClient *http.Client
}

// NewHealthService creates a HealthService implementation using Cloudflare.
func NewHealthService() domain.HealthService {
	return &cloudflareHealthService{
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

// CheckDomain implements domain.HealthService interface.
func (h *cloudflareHealthService) CheckDomain(ctx context.Context, domainURL string) (*models.DomainHealth, error) {
	parsedURL, err := url.Parse(domainURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	health := &models.DomainHealth{
		DNSResolved: false,
		SSLValid:    false,
		StatusCode:  0,
		LatencyMS:   0,
	}

	start := time.Now()
	resp, err := h.httpClient.Get(domainURL)
	latency := time.Since(start)
	health.LatencyMS = latency.Milliseconds()

	if err != nil {
		health.Message = fmt.Sprintf("DNS/Connection error: %v", err)
		return health, nil
	}
	defer resp.Body.Close()

	health.DNSResolved = true
	health.StatusCode = resp.StatusCode

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
