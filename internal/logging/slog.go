package logging

import (
	"context"
	"log/slog"
	"os"

	"hexleturlshort/internal/config"
)

// NewSlogLogger creates configured slog logger
func NewSlogLogger(env config.Environment) *slog.Logger {
	switch env {
	case config.EnvProduction:
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	default:
		return slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
}

type loggerKey struct{}

// ContextWithLogger adds logger to context
func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// FromContext takes slog.Logger from context
func FromContext(ctx context.Context, fallback *slog.Logger) *slog.Logger {
	l, ok := ctx.Value(loggerKey{}).(*slog.Logger)
	if !ok {
		return fallback
	}

	return l
}
