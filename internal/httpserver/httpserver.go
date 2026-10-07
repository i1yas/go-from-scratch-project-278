package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"

	"hexleturlshort/internal/config"
	"hexleturlshort/internal/httpserver/middleware"
)

// Server container http.Server and gin router
type Server struct {
	srv    *http.Server
	router *gin.Engine
}

// NewRouter generates configured gin router
func NewRouter(env config.Environment, logger *slog.Logger) *gin.Engine {
	setGinMode(env)

	r := gin.New()

	r.TrustedPlatform = gin.PlatformCloudflare

	r.Use(
		gin.Recovery(),
		sentrygin.New(sentrygin.Options{
			Repanic:         true,
			WaitForDelivery: false,
			Timeout:         5 * time.Second,
		}),
		middleware.RequestID(),
		middleware.Slog(logger),
		middleware.Timeout(),
		middleware.CORS(env),
	)

	return r
}

// NewServer creates Server based on gin router
func NewServer(router *gin.Engine, cfg config.HTTP) Server {
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: router.Handler(),
	}

	return Server{
		srv:    srv,
		router: router,
	}
}

// Run runs http.Server in Server
func (s *Server) Run() error {
	err := s.srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

// Shutdown shutdowns http.Server in Server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func setGinMode(env config.Environment) {
	switch env {
	case config.EnvProduction:
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}
