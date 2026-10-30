package application

import (
	"context"

	"hexleturlshort/internal/links"
)

// GetVisitsParams is params for GetVisits use-case
type GetVisitsParams struct {
	Range Range
	Sort  SortOrder
}

// VisitsResult result contains visits and total count in store
type VisitsResult struct {
	Items []links.Visit
	Total int64
}

var getVisitsSupportedSortField = map[string]struct{}{
	"id":         {},
	"link_id":    {},
	"ip":         {},
	"referer":    {},
	"user_agent": {},
	"created_at": {},
	"status":     {},
}

// GetVisits is use-case for listing visits
func (s *Service) GetVisits(ctx context.Context, params GetVisitsParams) (VisitsResult, error) {
	_, fieldIsSupported := getVisitsSupportedSortField[params.Sort.SortBy]
	if !fieldIsSupported {
		return VisitsResult{}, ErrUnsupportedSortOrder
	}

	totalCount, err := s.visits.GetVisitsTotalCount(ctx)
	if err != nil {
		return VisitsResult{}, err
	}

	if totalCount == 0 {
		return VisitsResult{
			Total: 0,
		}, nil
	}

	visits, err := s.visits.GetVisits(ctx, params)
	if err != nil {
		return VisitsResult{}, err
	}

	return VisitsResult{
		Items: visits,
		Total: totalCount,
	}, nil
}
