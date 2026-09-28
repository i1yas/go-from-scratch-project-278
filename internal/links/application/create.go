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
	originalURL, err := links.NewURL(params.OriginalURL)
	if err != nil {
		return links.Link{}, err
	}

	if params.ShortCode != nil {
		code, err := links.NewShortCode(*params.ShortCode)
		if err != nil {
			return links.Link{}, err
		}

		link := links.Link{
			OriginalURL: originalURL,
			ShortCode:   code,
		}

		return s.links.CreateLink(ctx, link)
	}

	const triesCount = 3

	for range triesCount {
		codeRaw, err := s.shortcodeGen.Generate()
		if err != nil {
			return links.Link{}, fmt.Errorf("%w: %w",
				ErrFailedToGenerateValidShortCode, err)
		}

		code, err := links.NewShortCode(codeRaw)
		if err != nil {
			return links.Link{}, fmt.Errorf("%w: %w",
				ErrFailedToGenerateValidShortCode, err)
		}

		link := links.Link{
			OriginalURL: originalURL,
			ShortCode:   code,
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
