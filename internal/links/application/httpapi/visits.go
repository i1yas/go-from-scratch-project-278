package httpapi

import (
	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LinksHandler handles HTTP requests related to links
type VisitsHandler struct {
	baseURL links.URL
	s       visitsService
}

// NewLinksHandler creates LinksHandler
func NewVisitsHandler(service visitsService, baseURL links.URL) *VisitsHandler {
	return &VisitsHandler{s: service, baseURL: baseURL}
}

func (h *VisitsHandler) GetVisits(ctx *gin.Context) {
	rang, err := parseRange(ctx.Query("range"), application.Range{From: 0, To: 4})
	if err != nil {
		handleError(ctx, err)
		return
	}

	sort, err := parseSort(ctx.Query("sort"), application.SortOrder{
		SortBy: "id",
		Order:  application.OrderASC,
	})
	if err != nil {
		handleError(ctx, err)
		return
	}

	result, err := h.s.GetVisits(ctx, application.GetVisitsParams{
		Range: rang,
		Sort:  sort,
	})
	if err != nil {
		handleError(ctx, err)
		return
	}

	response := make([]visitResponse, len(result.Items))

	for i, visit := range result.Items {
		response[i] = convertToVisitResponse(visit)
	}

	ctx.Header("Content-Range", formatContentRangeHeader(contentRangeParams{
		itemName:   "visits",
		itemsRange: rang,
		itemsCount: len(result.Items),
		totalCount: result.Total,
	}))
	ctx.JSON(http.StatusOK, response)
}

const timestampResponseFormat = "2006-01-02T15:04:05Z"

type visitResponse struct {
	ID        int64  `json:"id"`
	LinkID    int64  `json:"link_id"`
	IP        string `json:"ip"`
	Referer   string `json:"referer"`
	UserAgent string `json:"user-agent"`
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
		CreatedAt: visit.CreatedAt.Format(timestampResponseFormat),
		Status:    visit.Status,
	}
}
