package application

import (
	"context"

	"hexleturlshort/internal/links"
)

// GetLinksParams contains params for get links use-case
type GetLinksParams struct {
	Range Range
	Sort  SortOrder
}

// LinksResult contains total count and link items
type LinksResult struct {
	Total int64
	Items []links.Link
}

// GetLinks is use-case for listing links
func (s *Service) GetLinks(ctx context.Context, params GetLinksParams) (LinksResult, error) {
	totalCount, err := s.links.GetLinksTotalCount(ctx)
	if err != nil {
		return LinksResult{}, err
	}

	if totalCount == 0 {
		return LinksResult{
			Total: 0,
		}, nil
	}

	mapLinksSortToStore(&params.Sort)

	items, err := s.links.GetLinks(ctx, params)
	if err != nil {
		return LinksResult{}, err
	}

	return LinksResult{
		Total: totalCount,
		Items: items,
	}, nil
}

// GetLinkByID is use-case for loading link by id
func (s *Service) GetLinkByID(ctx context.Context, id int64) (links.Link, error) {
	return s.links.GetLinkByID(ctx, id)
}

func mapLinksSortToStore(sort *SortOrder) {
	if sort == nil {
		return
	}

	switch sort.sortBy {
	case "short_url", "short_name":
		sort.sortBy = "shortcode"
	}
}
