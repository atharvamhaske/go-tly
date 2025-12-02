package analytics

import (
	"net/http"
	"time"

	"github.com/atharvamhaske/go-tly/internal/domain/models"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, validate origin
	},
}

// Realtime handles WebSocket connections for real-time analytics updates.
type Realtime struct {
	queries  *Queries
	clients  map[*websocket.Conn]bool
	broadcast chan *models.AnalyticsSummary
}

// NewRealtime creates a new realtime handler.
func NewRealtime(queries *Queries) *Realtime {
	return &Realtime{
		queries:   queries,
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan *models.AnalyticsSummary, 256),
	}
}

// HandleWebSocket handles WebSocket connections for real-time updates.
func (r *Realtime) HandleWebSocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	r.clients[ws] = true
	defer delete(r.clients, ws)

	// Read short_key from query param
	shortKey := c.QueryParam("short_key")
	if shortKey == "" {
		return c.JSON(400, echo.Map{"error": "short_key required"})
	}

	// Send initial data
	summary, err := r.queries.GetAggregated(c.Request().Context(), shortKey)
	if err == nil {
		_ = ws.WriteJSON(summary)
	}

	// Update loop: poll every 5 seconds and send updates
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			summary, err := r.queries.GetAggregated(c.Request().Context(), shortKey)
			if err == nil {
				if err := ws.WriteJSON(summary); err != nil {
					return err
				}
			}

		case summary := <-r.broadcast:
			if err := ws.WriteJSON(summary); err != nil {
				return err
			}
		}
	}
}

// Broadcast sends an update to all connected clients.
func (r *Realtime) Broadcast(summary *models.AnalyticsSummary) {
	for client := range r.clients {
		_ = client.WriteJSON(summary)
	}
}

