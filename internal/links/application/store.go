package application

import (
	"context"

	"hexleturlshort/internal/links"
)

type linkStore interface {
	GetLinkByID(ctx context.Context, id int64) (links.Link, error)
	GetLinkByCode(ctx context.Context, code links.ShortCode) (links.Link, error)
	GetLinks(ctx context.Context, params GetLinksParams) ([]links.Link, error)
	GetLinksTotalCount(ctx context.Context) (int64, error)
	CreateLink(ctx context.Context, link links.Link) (links.Link, error)
	UpdateLink(ctx context.Context, link links.Link) (links.Link, error)
	DeleteLink(ctx context.Context, id int64) error
}

type visitsStore interface {
	GetVisits(ctx context.Context, params GetVisitsParams) ([]links.Visit, error)
	GetVisitsTotalCount(ctx context.Context) (int64, error)
	CreateVisit(ctx context.Context, visit links.Visit) error
}
