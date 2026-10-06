package httpserver

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"hexleturlshort/internal/config"
)

// Server container http.Server and gin router
type Server struct {
	srv    *http.Server
	router *gin.Engine
}

// NewRouter generates configured gin router
func NewRouter() *gin.Engine {
	r := gin.Default()

	r.TrustedPlatform = gin.PlatformCloudflare

	r.Use(gin.Recovery())

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

// TODO: configure properly
