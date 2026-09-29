package httpapi

import (
	"context"
	"errors"
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
			ctx.Status(400)
			return
		}

		result[i] = respItem
	}

	ctx.JSON(http.StatusOK, result)
}

// GetLinkByID handles request for getting link by id
func (h *LinksHandler) GetLinkByID(ctx *gin.Context) {
	idRaw := ctx.Param("id")

	id, err := strconv.ParseInt(idRaw, 10, 64)
	if err != nil {
		ctx.JSON(400, errorResponse{
			Message: "id should be number",
		})

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
	Message string `json:"message"`
}

func handleError(ctx *gin.Context, err error) {
	if errors.Is(err, application.ErrShortCodeConflict) {
		ctx.JSON(400, errorResponse{
			Message: "short name exist",
		})

		return
	}

	if errors.Is(err, application.ErrLinkNotFound) {
		ctx.JSON(404, errorResponse{
			Message: "link not found",
		})

		return
	}

	if errors.Is(err, links.ErrInvlalidShortCode) {
		ctx.JSON(404, errorResponse{
			Message: "invalid short name",
		})

		return
	}

	if errors.Is(err, links.ErrInvlalidURL) {
		ctx.JSON(404, errorResponse{
			Message: "invalid url",
		})

		return
	}

	ctx.Status(500)
}
