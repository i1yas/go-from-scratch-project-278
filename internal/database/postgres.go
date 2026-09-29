package database

import (
	"context"
	"database/sql"
	"time"

	// register pgx driver
	_ "github.com/jackc/pgx/v5/stdlib"

	"hexleturlshort/internal/config"
)

// OpenPostgres creates database handle and check connection
func OpenPostgres(ctx context.Context, cfg config.Database) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.URL)
	if err != nil {
		return nil, err
	}

	ctxPing, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err = db.PingContext(ctxPing)
	if err != nil {
		_ = db.Close() //nolint:errcheck // returns error when context canceled by test (t.Context)
		return nil, err
	}

	return db, nil
}

// TODO: configure db properly
