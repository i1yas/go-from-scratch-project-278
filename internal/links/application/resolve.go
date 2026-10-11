package application

import (
	"context"
	"log/slog"

	"hexleturlshort/internal/links"
)

// ResolveLinkParams contains information about link visit
type ResolveLinkParams struct {
	Code      string
	IP        string
	Referer   string
	UserAgent string
	Status    int
}

// ResolveLink is use-case for resolving short url
func (s *Service) ResolveLink(ctx context.Context, params ResolveLinkParams) (links.URL, error) {
	code, err := links.NewShortCode(params.Code)
	if err != nil {
		return links.URL(""), err
	}

	link, err := s.links.GetLinkByCode(ctx, code)
	if err != nil {
		return links.URL(""), err
	}

	visit, err := links.NewVisit(
		link.ID,
		params.IP,
		params.Referer,
		params.UserAgent,
		params.Status,
	)
	if err != nil {
		s.logger.Warn("failed to construct visit",
			slog.String("error", err.Error()),
		)
	}

	err = s.visits.CreateVisit(ctx, visit)
	if err != nil {
		s.logger.Error("failed to record visit",
			slog.String("error", err.Error()),
		)
	}

	return link.OriginalURL, nil
}
