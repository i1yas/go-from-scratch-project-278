package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Slog adds request logging with slog logger
func Slog(logger *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		logger.Info("request",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.Request.URL.Path),
			slog.Int("status", ctx.Writer.Status()),
			slog.Duration("latency", time.Since(start)),
		)
	}
}
