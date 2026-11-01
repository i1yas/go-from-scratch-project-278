package httpapi

import (
	"errors"
	"log/slog"
	"net/http"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
)

var (
	ErrInvalidID     = errors.New("invalid id")
	ErrInvalidBody   = errors.New("invalid body")
	ErrInvalidQuery  = errors.New("invalid query")
	ErrRangeIsTooBig = errors.New("range is too big")
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
		originalErr: ErrInvalidBody,
		status:      http.StatusBadRequest,
		Error:       "invalid request",
	},
	{
		originalErr: application.ErrInvalidRange,
		status:      http.StatusBadRequest,
		Error:       "invalid range",
	},
	{
		originalErr: ErrRangeIsTooBig,
		status:      http.StatusBadRequest,
		Error:       "range is too big",
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
		status:      http.StatusUnprocessableEntity,
		Errors: map[string]string{
			"short_name": "short name already in use",
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

func handleError(logger *slog.Logger, c *gin.Context, err error) {
	ok := errors.As(err, new(*http.MaxBytesError))
	if ok {
		c.Status(http.StatusRequestEntityTooLarge)
		return
	}

	var validationErrors validator.ValidationErrors

	if ok := errors.As(err, &validationErrors); ok {
		c.JSON(http.StatusUnprocessableEntity, errorResponse{
			Errors: validationErrorsMap(validationErrors),
		})

		return
	}

	for _, internalErr := range internalErrors {
		if errors.Is(err, internalErr) {
			handleInternalError(logger, c, err)
			return
		}
	}

	for _, mappedErr := range errorsMappings {
		if errors.Is(err, mappedErr.originalErr) {
			c.JSON(mappedErr.status, mappedErr)

			return
		}
	}

	handleInternalError(logger, c, err)
}

func handleInternalError(logger *slog.Logger, c *gin.Context, err error) {
	c.Status(500)
	logger.Error("internal error",
		slog.String("request_id", requestid.Get(c)),
		slog.String("error", err.Error()),
	)
	reportError(c, err)
}

func validationErrorsMap(errs validator.ValidationErrors) map[string]string {
	errsMap := make(map[string]string, len(errs))

	for _, vErr := range errs {
		errsMap[vErr.Field()] = vErr.Error()
	}

	return errsMap
}

func reportError(c *gin.Context, err error) {
	hub := sentrygin.GetHubFromContext(c)
	if hub != nil {
		hub.CaptureException(err)
	}
}
