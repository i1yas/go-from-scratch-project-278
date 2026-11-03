package postgres

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application/postgres/sqlcgen"
)

// migrationsDir returns path to migrations relative to this file
func migrationsDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("failed to build path to migrations")
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

func linksIDs(linkItems []links.Link) []int64 {
	ids := make([]int64, len(linkItems))

	for i, link := range linkItems {
		ids[i] = link.ID
	}

	return ids
}

func seedDB(t *testing.T, db sqlcgen.DBTX, name string) {
	path := filepath.Join("testdata", "seeds", name+".sql")

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	query := string(data)

	result, err := db.ExecContext(t.Context(), query)
	require.NoError(t, err)

	rows, err := result.RowsAffected()
	require.NoError(t, err)

	require.NotZero(t, rows)
}

func timeFromISO(t *testing.T, s string) time.Time {
	result, err := time.Parse(time.RFC3339, s)
	require.NoError(t, err)

	return result
}

func requireEqualVisit(t *testing.T, expected, actual links.Visit) {
	t.Helper()

	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.LinkID, actual.LinkID)
	require.Equal(t, expected.IP, actual.IP)
	require.Equal(t, expected.Referer, actual.Referer)
	require.Equal(t, expected.UserAgent, actual.UserAgent)
	require.Equal(t, expected.Status, actual.Status)
	require.Equal(t, expected.CreatedAt.UTC(), actual.CreatedAt.UTC())
}

func visitsIDs(visits []links.Visit) []int64 {
	ids := make([]int64, len(visits))

	for i, visit := range visits {
		ids[i] = visit.ID
	}

	return ids
}
