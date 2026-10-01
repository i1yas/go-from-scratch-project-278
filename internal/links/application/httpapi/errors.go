package httpapi

import (
	"errors"
	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
	"net/http"

	"github.com/gin-gonic/gin"
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
		Error:       "shortname already exist",
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

func handleError(ctx *gin.Context, err error) {
	for _, mappedErr := range errorsMappings {
		if errors.Is(err, mappedErr.originalErr) {
			ctx.JSON(mappedErr.status, mappedErr)

			return
		}
	}

	var linkErr *links.LinkError
	if errors.As(err, &linkErr) {
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse{
			Errors: linkErr.Fields,
		})

		return
	}

	ctx.Status(500)
}
