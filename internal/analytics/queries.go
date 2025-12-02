package analytics

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/atharvamhaske/go-tly/internal/domain/models"
)

// Queries provides aggregated analytics queries.
type Queries struct {
	conn driver.Conn
}

// NewQueries creates a new queries instance.
func NewQueries(conn driver.Conn) *Queries {
	return &Queries{conn: conn}
}

// GetAggregated returns aggregated analytics for a short key.
func (q *Queries) GetAggregated(ctx context.Context, shortKey string) (*models.AnalyticsSummary, error) {
	summary := &models.AnalyticsSummary{
		Countries: make(map[string]int64),
		Browsers:  make(map[string]int64),
		Devices:   make(map[string]int64),
	}

	// Total clicks
	var totalClicks int64
	if err := q.conn.QueryRow(ctx, `
		SELECT count() 
		FROM click_events 
		WHERE short_key = ?
	`, shortKey).Scan(&totalClicks); err != nil {
		return nil, fmt.Errorf("total clicks: %w", err)
	}
	summary.TotalClicks = totalClicks

	// Countries
	rows, err := q.conn.Query(ctx, `
		SELECT country, count() as cnt
		FROM click_events
		WHERE short_key = ?
		GROUP BY country
	`, shortKey)
	if err != nil {
		return nil, fmt.Errorf("countries query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var country string
		var cnt int64
		if err := rows.Scan(&country, &cnt); err != nil {
			continue
		}
		summary.Countries[country] = cnt
	}

	// Browsers
	rows, err = q.conn.Query(ctx, `
		SELECT browser, count() as cnt
		FROM click_events
		WHERE short_key = ?
		GROUP BY browser
	`, shortKey)
	if err != nil {
		return nil, fmt.Errorf("browsers query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var browser string
		var cnt int64
		if err := rows.Scan(&browser, &cnt); err != nil {
			continue
		}
		summary.Browsers[browser] = cnt
	}

	// Devices
	rows, err = q.conn.Query(ctx, `
		SELECT device, count() as cnt
		FROM click_events
		WHERE short_key = ?
		GROUP BY device
	`, shortKey)
	if err != nil {
		return nil, fmt.Errorf("devices query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var device string
		var cnt int64
		if err := rows.Scan(&device, &cnt); err != nil {
			continue
		}
		summary.Devices[device] = cnt
	}

	return summary, nil
}

