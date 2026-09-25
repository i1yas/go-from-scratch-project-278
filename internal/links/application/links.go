package application

import (
	"context"

	"hexleturlshort/internal/links"
)

// GetLinks is use-case for listing links
func (s *Service) GetLinks(ctx context.Context) ([]links.Link, error) {
	return s.links.GetLinks(ctx)
}

// GetLinkByID is use-case for loading link by id
func (s *Service) GetLinkByID(ctx context.Context, id int64) (links.Link, error) {
	return s.links.GetLinkByID(ctx, id)
}
