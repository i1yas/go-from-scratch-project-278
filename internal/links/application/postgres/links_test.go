package postgres

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
	"hexleturlshort/internal/links/application/postgres/sqlcgen"
)

func TestCreateAndGetLink(t *testing.T) {
	db := setupTestDB(t)

	cases := []struct {
		name  string
		input string
	}{
		{
			name:  "ascii",
			input: "test",
		},
		{
			name:  "unicode",
			input: "тест",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
				linksStore := NewLinksStore(tx)

				link := createValidLink(t, tc.input)

				createdLink, err := linksStore.CreateLink(ctx, link)
				require.NoError(t, err)
				require.Equal(t, link.OriginalURL, createdLink.OriginalURL)
				require.Equal(t, link.ShortCode, createdLink.ShortCode)

				loadedLinkByID, err := linksStore.GetLinkByID(ctx, createdLink.ID)
				require.NoError(t, err)
				require.Equal(t, createdLink, loadedLinkByID)

				loadedLinkByCode, err := linksStore.GetLinkByCode(ctx, createdLink.ShortCode)
				require.NoError(t, err)
				require.Equal(t, createdLink, loadedLinkByCode)
			})
		})
	}
}

func TestGetMultipleLinks(t *testing.T) {
	db := setupTestDB(t)

	withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
		linksStore := NewLinksStore(tx)

		link1 := createValidLink(t, "test-1")
		link2 := createValidLink(t, "test-2")
		link3 := createValidLink(t, "test-3")

		_, err := linksStore.CreateLink(ctx, link1)
		require.NoError(t, err)
		_, err = linksStore.CreateLink(ctx, link2)
		require.NoError(t, err)
		_, err = linksStore.CreateLink(ctx, link3)
		require.NoError(t, err)

		loadedLinks, err := linksStore.GetLinks(ctx)
		require.NoError(t, err)
		require.Equal(t, 3, len(loadedLinks))
	})
}

func TestUpdateLink(t *testing.T) {
	db := setupTestDB(t)

	t.Run("create and update", func(t *testing.T) {
		withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
			linksStore := NewLinksStore(tx)

			link := createValidLink(t, "test")

			createdLink, err := linksStore.CreateLink(ctx, link)
			require.NoError(t, err)

			update := createValidLink(t, "test-2")
			update.ID = createdLink.ID

			updatedLink, err := linksStore.UpdateLink(ctx, update)
			require.NoError(t, err)
			require.Equal(t, update.ShortCode, updatedLink.ShortCode)
			require.Equal(t, update.OriginalURL, updatedLink.OriginalURL)
		})
	})

	t.Run("error link not found", func(t *testing.T) {
		withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
			linksStore := NewLinksStore(tx)

			link := createValidLink(t, "test")

			_, err := linksStore.UpdateLink(ctx, link)
			require.ErrorIs(t, err, application.ErrLinkNotFound)
		})
	})

	t.Run("error on shortcode change, conflict", func(t *testing.T) {
		withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
			linksStore := NewLinksStore(tx)

			link1 := createValidLink(t, "test-1")
			link2 := createValidLink(t, "test-2")

			_, err := linksStore.CreateLink(ctx, link1)
			require.NoError(t, err)

			createdLink2, err := linksStore.CreateLink(ctx, link2)
			require.NoError(t, err)

			update := createdLink2
			update.ShortCode = link1.ShortCode

			_, err = linksStore.UpdateLink(ctx, update)
			require.ErrorIs(t, err, application.ErrShortCodeConflict)
		})
	})
}

func TestDeleteLink(t *testing.T) {
	db := setupTestDB(t)

	t.Run("create and delete", func(t *testing.T) {
		withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
			linksStore := NewLinksStore(tx)

			link := createValidLink(t, "test")

			createdLink, err := linksStore.CreateLink(ctx, link)
			require.NoError(t, err)

			err = linksStore.DeleteLink(ctx, createdLink.ID)
			require.NoError(t, err)
		})
	})

	t.Run("error link not found", func(t *testing.T) {
		withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
			linksStore := NewLinksStore(tx)

			err := linksStore.DeleteLink(ctx, 1)
			require.ErrorIs(t, err, application.ErrLinkNotFound)
		})
	})
}

func TestErrorLinkNotFound(t *testing.T) {
	db := setupTestDB(t)

	withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
		linksStore := NewLinksStore(tx)

		_, err := linksStore.GetLinkByID(ctx, 1)
		require.ErrorIs(t, err, application.ErrLinkNotFound)

		code, err := links.NewShortCode("does-not-exist")
		require.NoError(t, err)

		_, err = linksStore.GetLinkByCode(ctx, code)
		require.ErrorIs(t, err, application.ErrLinkNotFound)
	})
}

func TestErrorShortCodeConflict(t *testing.T) {
	db := setupTestDB(t)

	withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
		linksStore := NewLinksStore(tx)

		link := createValidLink(t, "test")

		_, err := linksStore.CreateLink(ctx, link)
		require.NoError(t, err)

		_, err = linksStore.CreateLink(ctx, link)
		require.ErrorIs(t, err, application.ErrShortCodeConflict)
	})
}

func TestErrorInvalidStoreValue(t *testing.T) {
	db := setupTestDB(t)

	withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
		linksStore := NewLinksStore(tx)

		link := createValidLink(t, "test")

		createdLink, err := linksStore.CreateLink(ctx, link)
		require.NoError(t, err)

		invalidCodeRaw := "x"
		_, err = links.NewShortCode(invalidCodeRaw)
		require.ErrorIs(t, err, links.ErrInvlalidShortCode)

		// NOTE: bypass store to write string directly
		_, err = linksStore.q.UpdateLink(ctx, sqlcgen.UpdateLinkParams{
			ID:          createdLink.ID,
			OriginalUrl: string(createdLink.OriginalURL),
			Shortcode:   invalidCodeRaw,
		})
		require.NoError(t, err)

		_, err = linksStore.GetLinkByID(ctx, createdLink.ID)
		require.ErrorIs(t, err, application.ErrInvalidStoreValue)
		require.ErrorIs(t, err, links.ErrInvlalidShortCode)
	})
}

func TestErrorInteralStore(t *testing.T) {
	db := setupTestDB(t)

	withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
		linksStore := NewLinksStore(tx)

		link := createValidLink(t, "test")

		err := tx.Rollback()
		require.NoError(t, err)

		_, err = linksStore.CreateLink(ctx, link)
		require.ErrorIs(t, err, application.ErrStoreInternal)
	})
}

var testDBContainer *postgres.PostgresContainer

func TestMain(m *testing.M) {
	m.Run()

	if testDBContainer != nil {
		err := testDBContainer.Terminate(context.Background())
		if err != nil {
			log.Fatal("failed to stop test db container")
		}
	}
}

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
