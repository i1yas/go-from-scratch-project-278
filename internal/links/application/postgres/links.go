package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
	"hexleturlshort/internal/links/application/postgres/sqlcgen"
)

const constraintUniqueShortCode = "unique_shortcode"

// LinksStore is postgres adapter for links store
type LinksStore struct {
	q sqlcgen.Querier
}

// NewLinksStore creates LinksStore
func NewLinksStore(db sqlcgen.DBTX) *LinksStore {
	querier := sqlcgen.New(db)
	return &LinksStore{q: querier}
}

// GetLinks loads links from store
func (s *LinksStore) GetLinks(ctx context.Context) ([]links.Link, error) {
	dbResult, err := s.q.GetLinks(ctx, sqlcgen.GetLinksParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", application.ErrStoreInternal, err)
	}

	result := make([]links.Link, len(dbResult))

	for i, dbLink := range dbResult {
		link, err := convertToLink(dbLink)
		if err != nil {
			return nil, err
		}

		result[i] = link
	}

	return result, nil
}

// GetLinksTotalCount returns total count of links
func (s *LinksStore) GetLinksTotalCount(ctx context.Context) (int64, error) {
	totalCount, err := s.q.GetLinksTotalCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", application.ErrStoreInternal, err)
	}

	return totalCount, nil
}

// GetLinkByID loads link from store by id
func (s *LinksStore) GetLinkByID(ctx context.Context, id int64) (links.Link, error) {
	dbLink, err := s.q.GetLinkByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return links.Link{}, application.ErrLinkNotFound
	}

	if err != nil {
		return links.Link{}, fmt.Errorf("%w: %w", application.ErrStoreInternal, err)
	}

	return convertToLink(dbLink)
}

// GetLinkByCode loads link from store by shortcode
func (s *LinksStore) GetLinkByCode(ctx context.Context, shortCode links.ShortCode) (links.Link, error) {
	dbLink, err := s.q.GetLinkByCode(ctx, string(shortCode))
	if errors.Is(err, sql.ErrNoRows) {
		return links.Link{}, application.ErrLinkNotFound
	}

	if err != nil {
		return links.Link{}, fmt.Errorf("%w: %w", application.ErrStoreInternal, err)
	}

	return convertToLink(dbLink)
}

// CreateLink adds link in store
func (s *LinksStore) CreateLink(ctx context.Context, link links.Link) (links.Link, error) {
	dbLink, err := s.q.CreateLink(ctx, sqlcgen.CreateLinkParams{
		OriginalUrl: string(link.OriginalURL),
		Shortcode:   string(link.ShortCode),
	})
	if isConstraintError(err, constraintUniqueShortCode) {
		return links.Link{}, application.ErrShortCodeConflict
	}

	if err != nil {
		return links.Link{}, fmt.Errorf("%w: %w", application.ErrStoreInternal, err)
	}

	return convertToLink(dbLink)
}

// UpdateLink updates link in store
func (s *LinksStore) UpdateLink(ctx context.Context, link links.Link) (links.Link, error) {
	dbLink, err := s.q.UpdateLink(ctx, sqlcgen.UpdateLinkParams{
		ID:          link.ID,
		OriginalUrl: string(link.OriginalURL),
		Shortcode:   string(link.ShortCode),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return links.Link{}, application.ErrLinkNotFound
	}

	if isConstraintError(err, constraintUniqueShortCode) {
		return links.Link{}, application.ErrShortCodeConflict
	}

	if err != nil {
		return links.Link{}, fmt.Errorf("%w: %w", application.ErrStoreInternal, err)
	}

	return convertToLink(dbLink)
}

// DeleteLink deletes link in store
func (s *LinksStore) DeleteLink(ctx context.Context, id int64) error {
	execRows, err := s.q.DeleteLink(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %w", application.ErrStoreInternal, err)
	}

	if execRows == 0 {
		return application.ErrLinkNotFound
	}

	return nil
}

func convertToLink(dbLink sqlcgen.Link) (links.Link, error) {
	link, err := links.NewLink(dbLink.OriginalUrl, dbLink.Shortcode)
	if err != nil {
		return links.Link{}, fmt.Errorf("%w: %w", application.ErrInvalidStoreValue, err)
	}

	link.ID = dbLink.ID

	return link, nil
}

func isConstraintError(err error, constraintName string) bool {
	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	if !ok {
		return false
	}

	return pgErr.ConstraintName == constraintName
}
