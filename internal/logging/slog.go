package logging

import (
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
