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

var getVisitsSupportedSort = map[string]string{
	"id ASC":          "id ASC",
	"id DESC":         "id DESC",
	"link_id ASC":     "link_id ASC, id ASC",
	"link_id DESC":    "link_id DESC, id ASC",
	"ip ASC":          "inet(ip) ASC, id ASC",
	"ip DESC":         "inet(ip) DESC, id ASC",
	"referer ASC":     "referer ASC NULLS FIRST, id ASC",
	"referer DESC":    "referer DESC NULLS LAST, id ASC",
	"user_agent ASC":  "user_agent ASC NULLS FIRST, id ASC",
	"user_agent DESC": "user_agent DESC NULLS LAST, id ASC",
	"created_at ASC":  "created_at ASC, id ASC",
	"created_at DESC": "created_at DESC, id ASC",
	"status ASC":      "status ASC, id ASC",
	"status DESC":     "status DESC, id ASC",
}

// GetVisits loads visits from store
func (s *VisitsStore) GetVisits(
	ctx context.Context,
	params application.GetVisitsParams,
) ([]links.Visit, error) {
	sortOrder, ok := getVisitsSupportedSort[params.Sort.String()]
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

	limit, offset := convertRangeToLimitOffset(params.Range)

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
			return nil, err
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

// CreateVisit creates visit in store
func (s *VisitsStore) CreateVisit(ctx context.Context, visit links.Visit) error {
	rowsAffected, err := s.q.CreateVisit(ctx, sqlcgen.CreateVisitParams{
		LinkID:    visit.LinkID,
		Ip:        addrToPgInet(visit.IP),
		Referer:   nullString(visit.Referer),
		UserAgent: nullString(visit.UserAgent),
		Status:    int16(visit.Status),
	})
	if err != nil {
		return fmt.Errorf("%w: %w", application.ErrStoreInternal, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: failed to insert visit", application.ErrStoreInternal)
	}

	return nil
}

func convertToVisit(dbVisit sqlcgen.Visit) (links.Visit, error) {
	visit, err := links.NewVisit(
		dbVisit.LinkID,
		dbVisit.Ip.IPNet.IP.String(),
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
