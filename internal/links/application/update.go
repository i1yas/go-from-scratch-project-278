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
	originalURL, err := links.NewURL(params.OriginalURL)
	if err != nil {
		return links.Link{}, err
	}

	code, err := links.NewShortCode(params.ShortCode)
	if err != nil {
		return links.Link{}, err
	}

	link := links.Link{
		ID:          params.ID,
		OriginalURL: originalURL,
		ShortCode:   code,
	}

	return s.links.UpdateLink(ctx, link)
}
