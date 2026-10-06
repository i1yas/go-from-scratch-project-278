package httpapi

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
)

// VisitsHandler handles HTTP requests related to links
type VisitsHandler struct {
	baseURL links.URL
	s       visitsService
}

// NewVisitsHandler creates VisitsHandler
func NewVisitsHandler(service visitsService, baseURL links.URL) *VisitsHandler {
	return &VisitsHandler{s: service, baseURL: baseURL}
}

// GetVisits lists recorded visits
func (h *VisitsHandler) GetVisits(c *gin.Context) {
	ctx := c.Request.Context()

	rang, err := parseRange(c.Query("range"), application.Range{From: 0, To: 4})
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

	result, err := h.s.GetVisits(ctx, application.GetVisitsParams{
		Range: rang,
		Sort:  sort,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	response := make([]visitResponse, len(result.Items))

	for i, visit := range result.Items {
		response[i] = convertToVisitResponse(visit)
	}

	c.Header("Content-Range", formatContentRangeHeader(contentRangeParams{
		itemName:   "visits",
		itemsRange: rang,
		itemsCount: len(result.Items),
		totalCount: result.Total,
	}))
	c.JSON(http.StatusOK, response)
}

type visitResponse struct {
	ID        int64  `json:"id"`
	LinkID    int64  `json:"link_id"`
	IP        string `json:"ip"`
	Referer   string `json:"reffer"` // NOTE: client expects `reffer` with typo
	UserAgent string `json:"user_agent"`
	CreatedAt string `json:"created_at"`
	Status    int    `json:"status"`
}

func convertToVisitResponse(visit links.Visit) visitResponse {
	return visitResponse{
		ID:        visit.ID,
		LinkID:    visit.LinkID,
		IP:        visit.IP.String(),
		Referer:   visit.Referer,
		UserAgent: visit.UserAgent,
		CreatedAt: visit.CreatedAt.UTC().Format(time.RFC3339),
		Status:    visit.Status,
	}
}
