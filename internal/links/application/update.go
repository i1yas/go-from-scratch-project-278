package application

import (
	"context"

	"hexleturlshort/internal/links"
)

// UpdateLinkParams is params for UpdateLink
type UpdateLinkParams struct {
	ID          int64
	OriginalURL string
	ShortCode   string
}

// UpdateLink is use-case for updating link
func (s *Service) UpdateLink(ctx context.Context, params UpdateLinkParams) (links.Link, error) {
	link, err := links.NewLinkWithID(params.ID, params.OriginalURL, params.ShortCode)
	if err != nil {
		return links.Link{}, err
	}

	return s.links.UpdateLink(ctx, link)
}
