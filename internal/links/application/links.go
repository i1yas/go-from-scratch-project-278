package application

import (
	"context"

	"hexleturlshort/internal/links"
)

// LinksResult contains total count and link items
type LinksResult struct {
	Total int64
	Items []links.Link
}

// GetLinks is use-case for listing links
func (s *Service) GetLinks(ctx context.Context) (LinksResult, error) {
	totalCount, err := s.links.GetLinksTotalCount(ctx)
	if err != nil {
		return LinksResult{}, err
	}

	if totalCount == 0 {
		return LinksResult{
			Total: 0,
		}, nil
	}

	items, err := s.links.GetLinks(ctx)
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
