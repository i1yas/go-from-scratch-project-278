package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
)

type linkService interface {
	GetLinks(ctx context.Context, params application.GetLinksParams) (application.LinksResult, error)
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
	linksRange, err := parseRange(ctx.Query("range"), application.Range{
		From: 0,
		To:   5,
	})
	if err != nil {
		handleError(ctx, err)
		return
	}

	// TODO: sort should be optional here, store should choose default one
	fallbackSort, err := application.NewSortOrder("id", application.OrderASC)
	if err != nil {
		handleError(ctx, err)
		return
	}

	sort, err := parseSort(ctx.Query("sort"), fallbackSort)
	if err != nil {
		handleError(ctx, err)
		return
	}

	linksResult, err := h.s.GetLinks(ctx, application.GetLinksParams{
		Range: linksRange,
		Sort:  sort,
	})
	if err != nil {
		handleError(ctx, err)
		return
	}

	linkItems := linksResult.Items
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
	ctx.Header("Content-Range", fmt.Sprintf("links 0-%d/%d", len(result), linksResult.Total))
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

// DeleteLink deletes link
func (h *LinksHandler) DeleteLink(ctx *gin.Context) {
	idRaw := ctx.Param("id")

	id, err := strconv.ParseInt(idRaw, 10, 64)
	if err != nil {
		handleError(ctx, fmt.Errorf("%w: %w", ErrInvalidID, err))
		return
	}

	if err := h.s.DeleteLink(ctx, id); err != nil {
		handleError(ctx, err)
		return
	}
}

// ResolveLink finds link by code and redirects to original url
func (h *LinksHandler) ResolveLink(ctx *gin.Context) {
	code := ctx.Param("code")

	originalURL, err := h.s.ResolveLink(ctx, code)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.Redirect(http.StatusFound, string(originalURL))
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

func parseRange(raw string, defaultRange application.Range) (application.Range, error) {
	if raw == "" {
		return defaultRange, nil
	}

	type rangeShape []int32

	var rangeItems rangeShape

	err := json.Unmarshal([]byte(raw), &rangeItems)
	if err != nil {
		return application.Range{}, ErrInvalidJSON
	}

	if len(rangeItems) < 2 {
		return defaultRange, nil
	}

	return application.NewRange(rangeItems[0], rangeItems[1])
}

func parseSort(raw string, defaultSort application.SortOrder) (application.SortOrder, error) {
	if raw == "" {
		return defaultSort, nil
	}

	type sortShape []string

	var sortItems sortShape

	err := json.Unmarshal([]byte(raw), &sortItems)
	if err != nil {
		return application.SortOrder{}, ErrInvalidJSON
	}

	if len(sortItems) < 2 {
		return defaultSort, nil
	}

	return application.NewSortOrder(sortItems[0], sortItems[1])
}
