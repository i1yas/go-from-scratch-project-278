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
	q  sqlcgen.Querier
	db sqlcgen.DBTX
}

// NewLinksStore creates LinksStore
func NewLinksStore(db sqlcgen.DBTX) *LinksStore {
	querier := sqlcgen.New(db)
	return &LinksStore{q: querier, db: db}
}

var getLinksSupportedSort = map[string]struct{}{
	"id ASC":            {},
	"id DESC":           {},
	"original_url ASC":  {},
	"original_url DESC": {},
	"shortcode ASC":     {},
	"shortcode DESC":    {},
}

// GetLinks loads links from store
func (s *LinksStore) GetLinks(ctx context.Context, params application.GetLinksParams) ([]links.Link, error) {
	sortOrder := params.Sort.String()

	_, ok := getLinksSupportedSort[sortOrder]
	if !ok {
		return nil, application.ErrUnsupportedSortOrder
	}

	query := fmt.Sprintf(`
		SELECT id, original_url, shortcode
		FROM links
		ORDER BY %s
		LIMIT $1 OFFSET $2;
	`, sortOrder)

	// NOTE: upper bound is inclusive
	limit := params.Range.To + 1 - params.Range.From

	offset := params.Range.From
	if params.Range.From == params.Range.To {
		limit = 0
	}

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", application.ErrStoreInternal, err)
	}

	defer rows.Close() //nolint:errcheck // Close error is checked via rows.Err()

	var items []links.Link

	for rows.Next() {
		var i sqlcgen.Link
		if err := rows.Scan(&i.ID, &i.OriginalUrl, &i.Shortcode); err != nil {
			return nil, fmt.Errorf("%w: %w", application.ErrStoreInternal, err)
		}

		link, err := convertToLink(i)
		if err != nil {
			return nil, err
		}

		items = append(items, link)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("%w: %w", application.ErrStoreInternal, err)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %w", application.ErrStoreInternal, err)
	}

	return items, nil
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
