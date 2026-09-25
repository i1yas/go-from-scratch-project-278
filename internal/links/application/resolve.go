package application

import (
	"context"

	"hexleturlshort/internal/links"
)

// ResolveLink is use-case for resolving short url
func (s *Service) ResolveLink(ctx context.Context, codeRaw string) (links.URL, error) {
	code, err := links.NewShortCode(codeRaw)
	if err != nil {
		return links.URL(""), err
	}

	link, err := s.links.GetLinkByCode(ctx, code)
	if err != nil {
		return links.URL(""), err
	}

	return link.OriginalURL, nil
}
