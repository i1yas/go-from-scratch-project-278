package httpapi

import (
	"errors"
	"net/http"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
)

var (
	ErrInvalidID    = errors.New("invalid id")
	ErrInvalidJSON  = errors.New("invalid json")
	ErrInvalidQuery = errors.New("invalid query")
)

type errorResponse struct {
	status      int
	originalErr error
	Error       string            `json:"error,omitempty"`
	Errors      map[string]string `json:"errors,omitempty"`
}

var errorsMappings = []errorResponse{
	{
		originalErr: ErrInvalidID,
		status:      http.StatusBadRequest,
		Error:       "invalid id",
	},
	{
		originalErr: ErrInvalidQuery,
		status:      http.StatusBadRequest,
		Error:       "invalid query",
	},
	{
		originalErr: ErrInvalidJSON,
		status:      http.StatusBadRequest,
		Error:       "invalid json",
	},
	{
		originalErr: application.ErrInvalidRange,
		status:      http.StatusBadRequest,
		Error:       "invalid range",
	},
	{
		originalErr: application.ErrInvalidSortOrder,
		status:      http.StatusBadRequest,
		Error:       "invalid sort order",
	},
	{
		originalErr: application.ErrUnsupportedSortOrder,
		status:      http.StatusBadRequest,
		Error:       "invalid sort order",
	},
	{
		originalErr: application.ErrShortCodeConflict,
		status:      http.StatusConflict,
		Errors: map[string]string{
			"short_name": "link with same short_name already exist",
		},
	},
	{
		originalErr: application.ErrLinkNotFound,
		status:      http.StatusNotFound,
		Error:       "link not found",
	},
	{
		originalErr: links.ErrInvlalidShortCode,
		status:      http.StatusUnprocessableEntity,
		Error:       "invalid shortname",
	},
	{
		originalErr: links.ErrInvlalidURL,
		status:      http.StatusUnprocessableEntity,
		Error:       "invalid url",
	},
}

var internalErrors = []error{
	application.ErrStoreInternal,
	application.ErrInvalidStoreValue,
	application.ErrShortCodeGeneratorInternal,
}

func handleError(c *gin.Context, err error) {
	for _, internalErr := range internalErrors {
		if errors.Is(err, internalErr) {
			c.Status(500)
			reportError(c, err)

			return
		}
	}

	for _, mappedErr := range errorsMappings {
		if errors.Is(err, mappedErr.originalErr) {
			c.JSON(mappedErr.status, mappedErr)

			return
		}
	}

	var linkErr *links.LinkError
	if errors.As(err, &linkErr) {
		c.JSON(http.StatusUnprocessableEntity, errorResponse{
			Errors: linkErr.Fields,
		})

		return
	}

	c.Status(500)
	reportError(c, err)
}

func reportError(c *gin.Context, err error) {
	hub := sentrygin.GetHubFromContext(c)
	if hub != nil {
		hub.CaptureException(err)
	}
}
