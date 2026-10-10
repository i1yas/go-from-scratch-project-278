package apiapp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"hexleturlshort/internal/config"
	"hexleturlshort/internal/database"
	"hexleturlshort/internal/errtrack"
	"hexleturlshort/internal/httpserver"
	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
	"hexleturlshort/internal/links/application/httpapi"
	"hexleturlshort/internal/links/application/postgres"
	"hexleturlshort/internal/links/application/shortcodegen"
)

var (
	ErrFailedToOpenDB            = errors.New("failed to open db")
	ErrFailedToReadConfigFromEnv = errors.New("failed to read config from environment")
	ErrInvalidBaseURL            = errors.New("invalid base url")
	ErrFailedToStartServer       = errors.New("failed to start server")
	ErrFailedToInitSentry        = errors.New("failed to init sentry client")
)

// App contains application components
type App struct {
	Cfg           config.Config
	DB            *sql.DB
	LinksStore    *postgres.LinksStore
	VisitsStore   *postgres.VisitsStore
	CodeGenerator *shortcodegen.Generator
	Service       *application.Service
	LinksHandler  *httpapi.LinksHandler
	Logger        *slog.Logger
	Router        *gin.Engine
	Server        httpserver.Server
}

// New wires application components and them in App struct
func New(ctx context.Context, cfg config.Config, logger *slog.Logger) (*App, error) {
	baseURL, err := links.NewURL(cfg.App.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidBaseURL, err)
	}

	db, err := database.OpenPostgres(ctx, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToOpenDB, err)
	}

	linksStore := postgres.NewLinksStore(db)
	visitsStore := postgres.NewVisitsStore(db)

	codeGenerator := shortcodegen.NewGenerator(shortcodegen.DefaultRandSource)

	service := application.NewService(
		linksStore,
		visitsStore,
		codeGenerator,
	)

	linksHandler := httpapi.NewLinksHandler(service, baseURL)
	visitsHandler := httpapi.NewVisitsHandler(service, baseURL)

	err = errtrack.InitSentry(errtrack.InitSentryParams{
		Env:    cfg.Env,
		Cfg:    cfg.Sentry,
		Logger: logger,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToInitSentry, err)
	}

	router := httpserver.NewRouter(cfg.Env, logger)
	server := httpserver.NewServer(router, cfg.HTTP)

	httpapi.RegisterRoutes(
		router,
		linksHandler,
		visitsHandler,
	)

	return &App{
		Cfg:           cfg,
		DB:            db,
		LinksStore:    linksStore,
		VisitsStore:   visitsStore,
		CodeGenerator: codeGenerator,
		Service:       service,
		LinksHandler:  linksHandler,
		Logger:        logger,
		Router:        router,
		Server:        server,
	}, nil
}

// Run runs application
func (app *App) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			time.Second*10,
		)
		defer cancel()

		err := app.Server.Shutdown(shutdownCtx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "falied to shutdown server: %v\n", err)
		}
	}()

	err := app.Server.Run()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToStartServer, err)
	}

	return nil
}

// Stop stops application
func (app *App) Stop(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(
		ctx,
		time.Second*10,
	)
	defer cancel()

	err := app.Server.Shutdown(shutdownCtx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "falied to shutdown server: %v\n", err)
	}
}
