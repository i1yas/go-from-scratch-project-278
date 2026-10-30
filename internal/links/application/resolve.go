package application

import (
	"context"
	"log/slog"

	"github.com/getsentry/sentry-go"

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

	s.recordVisit(ctx, link.ID, params)

	return link.OriginalURL, nil
}

func (s *Service) recordVisit(ctx context.Context, linkID int64, params ResolveLinkParams) {
	hub := sentry.GetHubFromContext(ctx)

	visit, err := links.NewVisit(
		linkID,
		params.IP,
		params.Referer,
		params.UserAgent,
		params.Status,
	)
	if err != nil {
		s.logger.Warn("failed to construct visit",
			slog.String("error", err.Error()),
		)
		reportError(hub, err)

		return
	}

	err = s.visits.CreateVisit(ctx, visit)
	if err != nil {
		s.logger.Error("failed to record visit",
			slog.String("error", err.Error()),
		)
		reportError(hub, err)

		return
	}
}
