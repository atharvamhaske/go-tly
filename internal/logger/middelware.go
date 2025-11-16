package logger

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GinLogger(l *Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		//store in the context
		ctx := c.Request.Context()
		ctx = contextWithRequestID(ctx, reqID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		latency := time.Since(start)

		l.Infof("incoming requests",
		"method", c.Request.Method,
		"path", c.FullPath(),
		"status", c.Writer.Status(),
		"latency_ms", latency.Milliseconds(),
		"client_ip", c.ClientIP(),
		"request_id", reqID,
		)
	}
}

func contextWithRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, "request_id", reqID)
}
