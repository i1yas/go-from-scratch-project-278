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

// GetVisits is use-case for listing visits
func (s *Service) GetVisits(ctx context.Context, params GetVisitsParams) (VisitsResult, error) {
	totalCount, err := s.visits.GetVisitsTotalCount(ctx)
	if err != nil {
		return VisitsResult{}, err
	}

	if totalCount == 0 {
		return VisitsResult{
			Total: 0,
		}, nil
	}

	mapVisitsSortToStore(&params.Sort)

	visits, err := s.visits.GetVisits(ctx, params)
	if err != nil {
		return VisitsResult{}, err
	}

	return VisitsResult{
		Items: visits,
		Total: totalCount,
	}, nil
}

func mapVisitsSortToStore(sort *SortOrder) {
	if sort == nil {
		return
	}

	if sort.SortBy == "reffer" {
		sort.SortBy = "referer"
	}
}
