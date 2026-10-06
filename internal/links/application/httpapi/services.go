package httpapi

import (
	"context"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
)

type linkService interface {
	GetLinks(ctx context.Context, params application.GetLinksParams) (application.LinksResult, error)
	GetLinkByID(ctx context.Context, id int64) (links.Link, error)
	CreateLink(ctx context.Context, params application.CreateLinkParams) (links.Link, error)
	UpdateLink(ctx context.Context, params application.UpdateLinkParams) (links.Link, error)
	DeleteLink(ctx context.Context, id int64) error
	ResolveLink(ctx context.Context, params application.ResolveLinkParams) (links.URL, error)
}

type visitsService interface {
	GetVisits(ctx context.Context, params application.GetVisitsParams) (application.VisitsResult, error)
}
