package postgres

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"hexleturlshort/internal/database/testutils"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	pg, err := testutils.SetupTestPostgres(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := pg.Migrate(migrationsDir()); err != nil {
		log.Fatal(err)
	}

	testDB = pg.DB

	m.Run()

	if err := pg.DB.Close(); err != nil {
		log.Fatal(err)
	}

	if err := pg.Container.Terminate(context.Background()); err != nil {
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
