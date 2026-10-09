package testutils

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// TestPostgres contains information about test postgres
type TestPostgres struct {
	Container *postgres.PostgresContainer
	DB        *sql.DB
	URL       string
}

// SetupTestPostgres starts postgres container and opens postgres
func SetupTestPostgres(ctx context.Context) (TestPostgres, error) {
	container, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("hexlet-shorturl-test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return TestPostgres{}, err
	}

	connString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return TestPostgres{}, err
	}

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return TestPostgres{}, err
	}

	if err := db.PingContext(ctx); err != nil {
		return TestPostgres{}, err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return TestPostgres{}, err
	}

	return TestPostgres{
		Container: container,
		DB:        db,
		URL:       connString,
	}, nil
}

// Migrate runs migrations on test postgres
func (tpg TestPostgres) Migrate(migrationsDir string) error {
	return goose.Up(tpg.DB, migrationsDir)
}
