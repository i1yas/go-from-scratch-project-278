package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	"hexleturlshort/internal/logging"
)

// Slog adds request logging with slog logger
func Slog(baseLogger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		logger := baseLogger.With("request_id", requestid.Get(c))
		ctx := logging.ContextWithLogger(c.Request.Context(), logger)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		logger.Info("request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", time.Since(start)),
		)
	}
}
