package middleware

import (
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// RequestID adds request_id
func RequestID() gin.HandlerFunc {
	return requestid.New(requestid.WithHandler(func(c *gin.Context, requestID string) {
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetTag("request_id", requestID)
			})
		}
	}))
}
