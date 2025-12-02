package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/atharvamhaske/go-tly/internal/domain/models"
)

// Pipeline batches click events and pushes them to ClickHouse.
type Pipeline struct {
	conn          driver.Conn
	batch         []*models.ClickEvents
	batchCh       chan *models.ClickEvents
	batchSize     int
	flushInterval time.Duration
}

// NewPipeline creates a new analytics pipeline.
func NewPipeline(dsn string, batchSize int, flushInterval time.Duration) (*Pipeline, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{dsn},
	})
	if err != nil {
		return nil, fmt.Errorf("clickhouse connect: %w", err)
	}

	p := &Pipeline{
		conn:          conn,
		batch:         make([]*models.ClickEvents, 0, batchSize),
		batchCh:       make(chan *models.ClickEvents, 1000),
		batchSize:     batchSize,
		flushInterval: flushInterval,
	}

	// Create table if not exists
	if err := p.createTable(context.Background()); err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}

	// Start batch processor
	go p.processBatches()

	return p, nil
}

// Push adds an event to the batch queue.
func (p *Pipeline) Push(ctx context.Context, evt *models.ClickEvents) error {
	select {
	case p.batchCh <- evt:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// processBatches collects events and flushes to ClickHouse periodically.
func (p *Pipeline) processBatches() {
	ticker := time.NewTicker(p.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case evt := <-p.batchCh:
			p.batch = append(p.batch, evt)
			if len(p.batch) >= p.batchSize {
				p.flush(context.Background())
			}

		case <-ticker.C:
			if len(p.batch) > 0 {
				p.flush(context.Background())
			}
		}
	}
}

// flush writes the current batch to ClickHouse.
func (p *Pipeline) flush(ctx context.Context) error {
	if len(p.batch) == 0 {
		return nil
	}

	batch, err := p.conn.PrepareBatch(ctx, "INSERT INTO click_events")
	if err != nil {
		return fmt.Errorf("prepare batch: %w", err)
	}

	for _, evt := range p.batch {
		if err := batch.Append(
			evt.ShortKey,
			evt.TimeStamp,
			evt.IP,
			evt.Country,
			evt.Browser,
			evt.Device,
			evt.Referrer,
			evt.UserAgent,
		); err != nil {
			return fmt.Errorf("append: %w", err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("send batch: %w", err)
	}

	// Clear batch
	p.batch = p.batch[:0]
	return nil
}

// createTable creates the click_events table in ClickHouse.
func (p *Pipeline) createTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS click_events (
			short_key String,
			timestamp DateTime,
			ip String,
			country String,
			browser String,
			device String,
			referrer String,
			user_agent String
		) ENGINE = MergeTree()
		ORDER BY (short_key, timestamp)
	`
	return p.conn.Exec(ctx, query)
}

