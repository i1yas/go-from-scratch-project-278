package apiapp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"hexleturlshort/internal/config"
	"hexleturlshort/internal/database"
	"hexleturlshort/internal/errtrack"
	"hexleturlshort/internal/httpserver"
	"hexleturlshort/internal/links/application"
	"hexleturlshort/internal/links/application/httpapi"
	"hexleturlshort/internal/links/application/postgres"
	"hexleturlshort/internal/links/application/shortcodegen"
)

var (
	ErrFailedToOpenDB            = errors.New("failed to open db")
	ErrFailedToReadConfigFromEnv = errors.New("failed to read config from environment")
	ErrFailedToCreateShortURLGen = errors.New("failed to create short url generator")
	ErrFailedToStartServer       = errors.New("failed to start server")
	ErrFailedToStopApp           = errors.New("failed to stop application")
	ErrFailedToInitSentry        = errors.New("failed to init sentry client")
)

// App contains application components
type App struct {
	Cfg    config.Config
	DB     *sql.DB
	Logger *slog.Logger
	Router *gin.Engine
	Server httpserver.Server
}

// New wires application components and them in App struct
func New(ctx context.Context, cfg config.Config, logger *slog.Logger) (*App, error) {
	shortURLGen, err := httpapi.NewShortURLGenerator(cfg.App.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateShortURLGen, err)
	}

	logger.Info("open postgres")

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
		logger,
	)

	linksHandler := httpapi.NewLinksHandler(service, shortURLGen, logger)
	visitsHandler := httpapi.NewVisitsHandler(service, logger)

	logger.Info("init senry")

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
		Cfg:    cfg,
		DB:     db,
		Logger: logger,
		Router: router,
		Server: server,
	}, nil
}

// Run runs application
func (app *App) Run(ctx context.Context) error {
	stopErr := make(chan error)

	go func() {
		defer close(stopErr)

		<-ctx.Done()

		app.Logger.Info("shutting down apiapp")

		stopErr <- app.Stop(context.Background())
	}()

	app.Logger.Info("starting server",
		slog.String("addr", app.Cfg.HTTP.Addr),
	)

	err := app.Server.Run()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%w: %w", ErrFailedToStartServer, err)
	}

	err = <-stopErr
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToStopApp, err)
	}

	return nil
}

// Stop stops application
func (app *App) Stop(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(
		ctx,
		time.Second*10,
	)
	defer cancel()

	app.Logger.Info("shutting down server")

	err := app.Server.Shutdown(shutdownCtx)
	if err != nil {
		return err
	}

	app.Logger.Info("closing postgres")

	err = app.DB.Close()
	if err != nil {
		return err
	}

	return nil
}
