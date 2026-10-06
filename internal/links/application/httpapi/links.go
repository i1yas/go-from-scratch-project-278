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
func (h *LinksHandler) GetLinks(c *gin.Context) {
	ctx := c.Request.Context()

	linksRange, err := parseRange(c.Query("range"), application.Range{
		From: 0,
		To:   4,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	sort, err := parseSort(c.Query("sort"), application.SortOrder{
		SortBy: "id",
		Order:  application.OrderASC,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	linksResult, err := h.s.GetLinks(ctx, application.GetLinksParams{
		Range: linksRange,
		Sort:  sort,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	linkItems := linksResult.Items
	result := make([]linkResponse, len(linkItems))

	for i, link := range linkItems {
		respItem, err := convertToLinkResponse(link, h.baseURL)
		if err != nil {
			handleError(c, err)
			return
		}

		result[i] = respItem
	}

	c.Header("Content-Range", formatContentRangeHeader(contentRangeParams{
		itemName:   "links",
		itemsRange: linksRange,
		itemsCount: len(result),
		totalCount: linksResult.Total,
	}))
	c.JSON(http.StatusOK, result)
}

type createLinkRequest struct {
	OriginalURL string  `json:"original_url"`
	ShortName   *string `json:"short_name,omitempty"`
}

// CreateLink creates link
func (h *LinksHandler) CreateLink(c *gin.Context) {
	ctx := c.Request.Context()

	var request createLinkRequest

	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		handleError(c, ErrInvalidJSON)
		return
	}

	link, err := h.s.CreateLink(ctx, application.CreateLinkParams{
		OriginalURL: request.OriginalURL,
		ShortCode:   request.ShortName,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	response, err := convertToLinkResponse(link, h.baseURL)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetLinkByID handles request for getting link by id
func (h *LinksHandler) GetLinkByID(c *gin.Context) {
	ctx := c.Request.Context()

	idRaw := c.Param("id")

	id, err := strconv.ParseInt(idRaw, 10, 64)
	if err != nil {
		handleError(c, fmt.Errorf("%w: %w", ErrInvalidID, err))
		return
	}

	link, err := h.s.GetLinkByID(ctx, id)
	if err != nil {
		handleError(c, err)
		return
	}

	resp, err := convertToLinkResponse(link, h.baseURL)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

type updateLinkRequest struct {
	OriginalURL string `json:"original_url"`
	ShortName   string `json:"short_name"`
}

// UpdateLink updates link
func (h *LinksHandler) UpdateLink(c *gin.Context) {
	ctx := c.Request.Context()

	idRaw := c.Param("id")

	id, err := strconv.ParseInt(idRaw, 10, 64)
	if err != nil {
		handleError(c, fmt.Errorf("%w: %w", ErrInvalidID, err))
		return
	}

	var request updateLinkRequest

	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		handleError(c, ErrInvalidJSON)
		return
	}

	link, err := h.s.UpdateLink(ctx, application.UpdateLinkParams{
		ID:          id,
		OriginalURL: request.OriginalURL,
		ShortCode:   request.ShortName,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	response, err := convertToLinkResponse(link, h.baseURL)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteLink deletes link
func (h *LinksHandler) DeleteLink(c *gin.Context) {
	ctx := c.Request.Context()

	idRaw := c.Param("id")

	id, err := strconv.ParseInt(idRaw, 10, 64)
	if err != nil {
		handleError(c, fmt.Errorf("%w: %w", ErrInvalidID, err))
		return
	}

	if err := h.s.DeleteLink(ctx, id); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ResolveLink finds link by code and redirects to original url
func (h *LinksHandler) ResolveLink(c *gin.Context) {
	ctx := c.Request.Context()

	code := c.Param("code")

	status := http.StatusFound

	originalURL, err := h.s.ResolveLink(ctx, application.ResolveLinkParams{
		Code:      code,
		IP:        c.ClientIP(),
		Referer:   c.Request.Referer(),
		UserAgent: c.Request.UserAgent(),
		Status:    status,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.Redirect(status, string(originalURL))
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
