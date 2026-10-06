package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"hexleturlshort/internal/config"
)

const frontendDevOrigin = "http://localhost:5173"

// CORS configures CORS policies
func CORS(env config.Environment) gin.HandlerFunc {
	switch env {
	case config.EnvProduction:
		// NOTE: no CORS allowed in production
		return func(*gin.Context) {}
	default:
		return cors.New(cors.Config{
			AllowOrigins:     []string{frontendDevOrigin},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Range", "Referer"},
			AllowCredentials: false,
			MaxAge:           12 * time.Hour,
		})
	}
}
