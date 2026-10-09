package postgres

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("hexlet-shorturl-test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatal(err)
	}

	connString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("pgx", connString)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	migrations := migrationsDir()
	if err := goose.Up(db, migrations); err != nil {
		log.Fatal(err)
	}

	testDB = db

	m.Run()

	if err := db.Close(); err != nil {
		log.Fatal(err)
	}

	if err := container.Terminate(context.Background()); err != nil {
		log.Fatal("failed to stop test db container")
	}
}

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	if testDB == nil {
		log.Fatal("missing test db")
	}

	return testDB
}
