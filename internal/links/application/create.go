package application

import (
	"context"
	"errors"
	"fmt"

	"hexleturlshort/internal/links"
)

// CreateLinkParams is params for CreateLink
type CreateLinkParams struct {
	OriginalURL string
	ShortCode   *string
}

// CreateLink is use-case for creating link entry
func (s *Service) CreateLink(ctx context.Context, params CreateLinkParams) (links.Link, error) {
	if params.ShortCode != nil {
		link, err := links.NewLink(params.OriginalURL, *params.ShortCode)
		if err != nil {
			return links.Link{}, err
		}

		return s.links.CreateLink(ctx, link)
	}

	const triesCount = 3

	for range triesCount {
		generatedCode, err := s.shortcodeGen.Generate()
		if err != nil {
			return links.Link{}, fmt.Errorf("%w: %w",
				ErrFailedToGenerateValidShortCode, err)
		}

		link, err := links.NewLink(params.OriginalURL, generatedCode)
		if err != nil {
			return links.Link{}, err
		}

		createdLink, err := s.links.CreateLink(ctx, link)
		if errors.Is(err, ErrShortCodeConflict) {
			continue
		}

		if err != nil {
			return links.Link{}, err
		}

		return createdLink, nil
	}

	return links.Link{}, ErrFailedToGenerateUniqueShortCode
}
