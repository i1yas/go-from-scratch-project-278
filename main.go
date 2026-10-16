package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"hexleturlshort/internal/apiapp"
	"hexleturlshort/internal/config"
	"hexleturlshort/internal/logging"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.ReadFromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read config: %v\n", err)
		os.Exit(1)
	}

	logger := logging.NewSlogLogger(cfg.Env)

	logger.Info("creating apiapp",
		slog.String("environment", string(cfg.Env)),
	)

	app, err := apiapp.New(ctx, cfg, logger)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	logger.Info("starting apiapp")

	err = app.Run(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
