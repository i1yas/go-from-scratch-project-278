package errtrack

import (
	"log/slog"

	"github.com/getsentry/sentry-go"

	"hexleturlshort/internal/config"
)

// InitSentryParams contains params for sentry client init
type InitSentryParams struct {
	Env    config.Environment
	Cfg    config.Sentry
	Logger *slog.Logger
}

// InitSentry inits sentry client
func InitSentry(params InitSentryParams) error {
	logger := params.Logger

	switch params.Env {
	case config.EnvProduction:
		if params.Cfg.DSN == "" {
			logger.Warn("sentry init: no DSN provided")
		}

		return sentry.Init(sentry.ClientOptions{
			Dsn: params.Cfg.DSN,
		})
	default:
		return sentry.Init(sentry.ClientOptions{
			Dsn:   params.Cfg.DSN,
			Debug: false,
			BeforeSend: func(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
				logger.Info("sentry event",
					slog.String("error", event.Message),
					slog.Any("exception", event.Exception),
					slog.Any("tags", event.Tags),
				)

				return nil
			},
		})
	}
}
