package httpapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
)

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
	linksRange, err := parseRange(ctx.Query("range"), application.Range{
		From: 0,
		To:   4,
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

	ctx.Header("Content-Range", formatContentRangeHeader(contentRangeParams{
		itemName:   "links",
		itemsRange: linksRange,
		itemsCount: len(result),
		totalCount: linksResult.Total,
	}))
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

	ctx.Status(http.StatusNoContent)
}

// ResolveLink finds link by code and redirects to original url
func (h *LinksHandler) ResolveLink(ctx *gin.Context) {
	code := ctx.Param("code")

	status := http.StatusFound

	originalURL, err := h.s.ResolveLink(ctx, application.ResolveLinkParams{
		Code:      code,
		IP:        ctx.ClientIP(),
		Referer:   ctx.Request.Referer(),
		UserAgent: ctx.Request.UserAgent(),
		Status:    status,
	})
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.Redirect(status, string(originalURL))
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

type contentRangeParams struct {
	itemName   string
	itemsRange application.Range
	itemsCount int
	totalCount int64
}

func formatContentRangeHeader(params contentRangeParams) string {
	r := params.itemsRange

	return fmt.Sprintf("%s %d-%d/%d",
		params.itemName,
		r.From,
		r.From+int32(params.itemsCount)-1,
		params.totalCount,
	)
}
