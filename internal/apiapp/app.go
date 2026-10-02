package apiapp

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"hexleturlshort/internal/config"
	"hexleturlshort/internal/database"
	"hexleturlshort/internal/httpserver"
	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
	"hexleturlshort/internal/links/application/httpapi"
	"hexleturlshort/internal/links/application/postgres"
	"hexleturlshort/internal/links/application/shortcodegen"
)

var (
	ErrFailedToReadConfigFromEnv = errors.New("failed to read config from environment")
	ErrInvalidBaseURL            = errors.New("invalid base url")
	ErrFailedToStartServer       = errors.New("failed to start server")
)

// Run combines components in one application and runs it
func Run(ctx context.Context) error {
	cfg, err := config.ReadFromEnv()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToReadConfigFromEnv, err)
	}

	db, err := database.OpenPostgres(ctx, cfg.Database)
	if err != nil {
		return err
	}

	linksStore := postgres.NewLinksStore(db)
	visitsStore := postgres.NewVisitsStore(db)

	codeGenerator := shortcodegen.NewGenerator(shortcodegen.DefaultRandSource)

	service := application.NewService(
		linksStore,
		visitsStore,
		codeGenerator,
	)

	baseURL, err := links.NewURL(cfg.App.BaseURL)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidBaseURL, err)
	}

	linksHandler := httpapi.NewLinksHandler(service, baseURL)

	router := httpserver.NewRouter()
	server := httpserver.NewServer(router, cfg.HTTP)

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			time.Second*10,
		)
		defer cancel()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "falied to shutdown server: %v\n", err)
		}
	}()

	httpapi.RegisterRoutes(
		router,
		linksHandler,
	)

	err = server.Run()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToStartServer, err)
	}

	return nil
}
