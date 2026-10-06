package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

const requestTimeout = 5 * time.Second

// Timeout adds request timeout
func Timeout() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
