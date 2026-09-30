package httpapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
)

type linkService interface {
	GetLinks(ctx context.Context) ([]links.Link, error)
	GetLinkByID(ctx context.Context, id int64) (links.Link, error)
	CreateLink(ctx context.Context, params application.CreateLinkParams) (links.Link, error)
	UpdateLink(ctx context.Context, params application.UpdateLinkParams) (links.Link, error)
	DeleteLink(ctx context.Context, id int64) error
	ResolveLink(ctx context.Context, code string) (links.URL, error)
}

// LinksHandler handles HTTP requests related to links
type LinksHandler struct {
	baseURL links.URL
	s       linkService
}

var (
	ErrInvalidID   = errors.New("invalid id")
	ErrInvalidJSON = errors.New("invalid json")
)

// NewLinksHandler creates LinksHandler
func NewLinksHandler(service linkService, baseURL links.URL) *LinksHandler {
	return &LinksHandler{s: service, baseURL: baseURL}
}

// GetLinks handles request for listing links
func (h *LinksHandler) GetLinks(ctx *gin.Context) {
	linkItems, err := h.s.GetLinks(ctx)
	if err != nil {
		handleError(ctx, err)
		return
	}

	result := make([]linkResponse, len(linkItems))

	for i, link := range linkItems {
		respItem, err := convertToLinkResponse(link, h.baseURL)
		if err != nil {
			handleError(ctx, err)
			return
		}

		result[i] = respItem
	}

	// TODO: implement properly
	ctx.Header("Content-Range", fmt.Sprintf("links 0-%d/%d", len(result), len(result)))
	ctx.JSON(http.StatusOK, result)
}

type createLinkRequest struct {
	OriginalURL string  `json:"original_url"`
	ShortName   *string `json:"short_name,omitempty"`
}

// CreateLink creates link
func (h *LinksHandler) CreateLink(ctx *gin.Context) {
	var request createLinkRequest

	if err := ctx.ShouldBindBodyWithJSON(&request); err != nil {
		handleError(ctx, ErrInvalidJSON)
		return
	}

	link, err := h.s.CreateLink(ctx, application.CreateLinkParams{
		OriginalURL: request.OriginalURL,
		ShortCode:   request.ShortName,
	})
	if err != nil {
		handleError(ctx, err)
		return
	}

	response, err := convertToLinkResponse(link, h.baseURL)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetLinkByID handles request for getting link by id
func (h *LinksHandler) GetLinkByID(ctx *gin.Context) {
	idRaw := ctx.Param("id")

	id, err := strconv.ParseInt(idRaw, 10, 64)
	if err != nil {
		handleError(ctx, fmt.Errorf("%w: %w", ErrInvalidID, err))
		return
	}

	link, err := h.s.GetLinkByID(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	resp, err := convertToLinkResponse(link, h.baseURL)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

type updateLinkRequest struct {
	OriginalURL string `json:"original_url"`
	ShortName   string `json:"short_name"`
}

// UpdateLink updates link
func (h *LinksHandler) UpdateLink(ctx *gin.Context) {
	idRaw := ctx.Param("id")

	id, err := strconv.ParseInt(idRaw, 10, 64)
	if err != nil {
		handleError(ctx, fmt.Errorf("%w: %w", ErrInvalidID, err))
		return
	}

	var request updateLinkRequest

	if err := ctx.ShouldBindBodyWithJSON(&request); err != nil {
		handleError(ctx, ErrInvalidJSON)
		return
	}

	link, err := h.s.UpdateLink(ctx, application.UpdateLinkParams{
		ID:          id,
		OriginalURL: request.OriginalURL,
		ShortCode:   request.ShortName,
	})
	if err != nil {
		handleError(ctx, err)
		return
	}

	response, err := convertToLinkResponse(link, h.baseURL)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

type linkResponse struct {
	ID          int64  `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortName   string `json:"short_name"`
	ShortURL    string `json:"short_url"`
}

func convertToLinkResponse(link links.Link, baseURL links.URL) (linkResponse, error) {
	shortURL, err := link.ShortURL(baseURL)
	if err != nil {
		return linkResponse{}, err
	}

	return linkResponse{
		ID:          link.ID,
		OriginalURL: string(link.OriginalURL),
		ShortName:   string(link.ShortCode),
		ShortURL:    string(shortURL),
	}, nil
}

type errorResponse struct {
	Error  string            `json:"error,omitempty"`
	Errors map[string]string `json:"errors,omitempty"`
}

func handleError(ctx *gin.Context, err error) {
	var linkErr *links.LinkError

	switch {
	case errors.Is(err, ErrInvalidID):
		ctx.JSON(http.StatusBadRequest, errorResponse{
			Error: "invalid id",
		})

	case errors.Is(err, ErrInvalidJSON):
		ctx.JSON(http.StatusBadRequest, errorResponse{
			Error: "invalid json",
		})

	case errors.Is(err, application.ErrShortCodeConflict):
		ctx.JSON(http.StatusConflict, errorResponse{
			Error: "shortname already exist",
		})

	case errors.Is(err, application.ErrLinkNotFound):
		ctx.JSON(http.StatusNotFound, errorResponse{
			Error: "link not found",
		})

	case errors.As(err, &linkErr):
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse{
			Errors: linkErr.Fields,
		})

	case errors.Is(err, links.ErrInvlalidShortCode):
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse{
			Error: "invalid shortname",
		})

	case errors.Is(err, links.ErrInvlalidURL):
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse{
			Error: "invalid url",
		})

	default:
		ctx.Status(500)
	}
}
