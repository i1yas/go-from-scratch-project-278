package postgres

import (
	"context"
	"database/sql"
	"hexleturlshort/internal/links"
	"log"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var testDBContainer *postgres.PostgresContainer

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	ctx := t.Context()

	if testDBContainer == nil {
		container, err := postgres.Run(
			ctx,
			"postgres:16-alpine",
			postgres.WithDatabase("hexlet-shorturl-test-db"),
			postgres.WithUsername("postgres"),
			postgres.WithPassword("postgres"),
			postgres.BasicWaitStrategies(),
		)
		require.NoError(t, err)

		testDBContainer = container
	}

	connString, err := testDBContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sql.Open("pgx", connString)
	require.NoError(t, err)

	require.NoError(t, db.PingContext(ctx))

	require.NoError(t, goose.SetDialect("postgres"))

	migrations := migrationsDir(t)
	require.NoError(t, goose.Up(db, migrations))

	t.Cleanup(func() {
		err := db.Close()
		require.NoError(t, err)
	})

	return db
}

func stopTestDBContainer() {
	if testDBContainer != nil {
		err := testDBContainer.Terminate(context.Background())
		if err != nil {
			log.Fatal("failed to stop test db container")
		}
	}
}

// migrationsDir returns path to migrations relative to this file
func migrationsDir(t *testing.T) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to build path to migrations")
	}

	return filepath.Join(
		filename,
		"..", "..", "..", "..", "..",
		"db", "migrations",
	)
}

func withTx(t *testing.T, db *sql.DB, fn func(ctx context.Context, tx *sql.Tx)) {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	t.Cleanup(cancel)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = tx.Rollback() //nolint:errcheck // returns error when context canceled by test (t.Context)
	})

	fn(ctx, tx)
}

func createValidLink(t *testing.T, codeRaw string) links.Link {
	t.Helper()

	originalURL, err := links.NewURL("https://test.com")
	require.NoError(t, err)

	code, err := links.NewShortCode(codeRaw)
	require.NoError(t, err)

	return links.Link{OriginalURL: originalURL, ShortCode: code}
}

func linksCodes(linkItems []links.Link) []links.ShortCode {
	codes := make([]links.ShortCode, len(linkItems))

	for i, link := range linkItems {
		codes[i] = link.ShortCode
	}

	return codes
}
