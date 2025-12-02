package analytics

import (
	"net/http"
	"time"

	"github.com/atharvamhaske/go-tly/internal/cloudflare"
	"github.com/atharvamhaske/go-tly/internal/domain/models"
	"github.com/labstack/echo/v4"
)

// Ingester handles incoming click event logs from HTTP requests.
// It can receive events from:
// 1. Direct HTTP POST (for testing/fallback)
// 2. Cloudflare Worker (with signature verification)
type Ingester struct {
	pipeline    *Pipeline
	logVerifier *cloudflare.WorkerLogVerifier
}

// NewIngester creates a new analytics ingester.
func NewIngester(pipeline *Pipeline, logVerifier *cloudflare.WorkerLogVerifier) *Ingester {
	return &Ingester{
		pipeline:    pipeline,
		logVerifier: logVerifier,
	}
}

// LogClick handles POST /api/analytics/click
// Receives click events and pushes them to the pipeline.
func (i *Ingester) LogClick(c echo.Context) error {
	// Check if this is a Worker-signed log
	timestamp := c.Request().Header.Get("X-Worker-Timestamp")
	signature := c.Request().Header.Get("X-Worker-Signature")
	body := c.Request().Header.Get("X-Worker-Body")

	var evt *models.ClickEvents
	var err error

	if timestamp != "" && signature != "" && body != "" {
		// Verify Worker log signature
		evt, err = i.logVerifier.VerifyLogEntry(timestamp, signature, body)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid signature"})
		}
	} else {
		// Direct HTTP (for testing/fallback)
		var directEvt models.ClickEvents
		if err := c.Bind(&directEvt); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
		}
		evt = &directEvt

		// Set timestamp if not provided
		if evt.TimeStamp.IsZero() {
			evt.TimeStamp = time.Now()
		}

		// Extract IP from request if not provided
		if evt.IP == "" {
			evt.IP = c.RealIP()
		}

		// Extract User-Agent if not provided
		if evt.UserAgent == "" {
			evt.UserAgent = c.Request().UserAgent()
		}
	}

	// Push to pipeline (async, non-blocking)
	go func() {
		_ = i.pipeline.Push(c.Request().Context(), evt)
	}()

	return c.JSON(http.StatusAccepted, echo.Map{"status": "accepted"})
}

