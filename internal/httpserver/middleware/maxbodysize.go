package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MaxBodySize adds limit for request body
func MaxBodySize(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		c.Next()
	}
}
