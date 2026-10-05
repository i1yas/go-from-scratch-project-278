package postgres

import (
	"context"
	"fmt"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
	"hexleturlshort/internal/links/application/postgres/sqlcgen"
)

// VisitsStore is postgres adapter for visits store
type VisitsStore struct {
	q  sqlcgen.Querier
	db sqlcgen.DBTX
}

// NewVisitsStore creates VisitsStore
func NewVisitsStore(db sqlcgen.DBTX) *VisitsStore {
	querier := sqlcgen.New(db)
	return &VisitsStore{q: querier, db: db}
}

var getVisitsSupportedSort = map[string]struct{}{
	"id ASC":  {},
	"id DESC": {},
}

// GetVisits loads visits from store
func (s *VisitsStore) GetVisits(
	ctx context.Context,
	params application.GetVisitsParams,
) ([]links.Visit, error) {
	sortOrder := params.Sort.String()

	_, ok := getVisitsSupportedSort[sortOrder]
	if !ok {
		return nil, application.ErrUnsupportedSortOrder
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			link_id,
			created_at,
			ip,
			referer,
			user_agent,
			status
		FROM visits
		ORDER BY %s
		LIMIT $1 OFFSET $2;
	`, sortOrder)

	// NOTE: upper bound is inclusive
	limit := params.Range.To + 1 - params.Range.From

	offset := max(0, params.Range.From)
	if params.Range.From == params.Range.To {
		limit = 0
	}

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", application.ErrStoreInternal, err)
	}

	defer rows.Close() //nolint:errcheck // Close error is checked via rows.Err()

	var items []links.Visit

	for rows.Next() {
		var i sqlcgen.Visit
		if err := rows.Scan(
			&i.ID,
			&i.LinkID,
			&i.CreatedAt,
			&i.Ip,
			&i.Referer,
			&i.UserAgent,
			&i.Status,
		); err != nil {
			return nil, fmt.Errorf("%w: %w", application.ErrStoreInternal, err)
		}

		visit, err := convertToVisit(i)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", application.ErrInvalidStoreValue, err)
		}

		items = append(items, visit)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("%w: %w", application.ErrStoreInternal, err)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %w", application.ErrStoreInternal, err)
	}

	return items, nil
}

// GetVisitsTotalCount returns total visits count in store
func (s *VisitsStore) GetVisitsTotalCount(ctx context.Context) (int64, error) {
	totalCount, err := s.q.GetVisitsTotalCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", application.ErrStoreInternal, err)
	}

	return totalCount, nil
}

func convertToVisit(dbVisit sqlcgen.Visit) (links.Visit, error) {
	visit, err := links.NewVisit(
		dbVisit.LinkID,
		dbVisit.Ip.IPNet.String(),
		dbVisit.Referer.String,
		dbVisit.UserAgent.String,
		int(dbVisit.Status),
	)
	if err != nil {
		return links.Visit{}, fmt.Errorf("%w: %w", application.ErrInvalidStoreValue, err)
	}

	visit.ID = dbVisit.ID
	visit.CreatedAt = dbVisit.CreatedAt

	return visit, nil
}
