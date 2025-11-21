package logger

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (l *Logger) EchoLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			res := c.Response()

			//genarate request ID
			reqID := req.Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = uuid.NewString()
			}
			c.Set("request_id", reqID)

			err := next(c)

			latency := time.Since(start)

			l.Infow("incoming request",
				"request_id", reqID,
				"method", req.Method,
				"path", req.URL.Path,
				"status", res.Status,
				"latency_ms", latency.Milliseconds(),
				"ip", c.RealIP(),
			)
			return err
		}
	}
}
