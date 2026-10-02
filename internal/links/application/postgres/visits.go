package postgres

import (
	"context"

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

// GetVisits loads visits from store
func (s *VisitsStore) GetVisits(
	context.Context,
	application.GetVisitsParams,
) ([]links.Visit, error) {
	// TODO: implement this
	return nil, nil
}

// GetVisitsTotalCount returns total visits count in store
func (s *VisitsStore) GetVisitsTotalCount(context.Context) (int64, error) {
	// TODO: implement this
	return 0, nil
}
