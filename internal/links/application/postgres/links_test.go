package postgres

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
	"hexleturlshort/internal/links/application/postgres/sqlcgen"
	"hexleturlshort/internal/links/application/testutils"
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

				link := testutils.Link(t, tc.input)

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

func TestGetLinksPagination(t *testing.T) {
	cases := []struct {
		name        string
		from        int
		to          int
		resultCount int
	}{
		{
			name:        "from 0 to 3",
			from:        0,
			to:          3,
			resultCount: 4,
		},
		{
			name:        "from 0 to 3, total 4",
			from:        0,
			to:          3,
			resultCount: 4,
		},
		{
			name:        "from 30 to 40, total 0",
			from:        30,
			to:          40,
			resultCount: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)

			withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
				seedDB(t, tx, "links")

				linksStore := NewLinksStore(tx)

				linksRange, err := application.NewRange(int32(tc.from), int32(tc.to))
				require.NoError(t, err)

				loadedLinks, err := linksStore.GetLinks(ctx, application.GetLinksParams{
					Range: linksRange,
					Sort:  testutils.SortOrder(t, "id", "ASC"),
				})

				require.NoError(t, err)
				require.Equal(t, tc.resultCount, len(loadedLinks))
			})
		})
	}
}

func TestGetLinksSorting(t *testing.T) {
	cases := []struct {
		name        string
		sortBy      string
		order       string
		wantIDOrder []int64
	}{
		{
			name:        "id ASC",
			sortBy:      "id",
			order:       application.OrderASC,
			wantIDOrder: []int64{1, 2, 3},
		},
		{
			name:        "id DESC",
			sortBy:      "id",
			order:       application.OrderDESC,
			wantIDOrder: []int64{10, 9, 8},
		},
		{
			name:        "shortcode ASC",
			sortBy:      "shortcode",
			order:       application.OrderASC,
			wantIDOrder: []int64{9, 10, 3},
		},
		{
			name:        "shortcode DESC",
			sortBy:      "shortcode",
			order:       application.OrderDESC,
			wantIDOrder: []int64{2, 5, 8},
		},
		{
			name:        "original_url ASC",
			sortBy:      "original_url",
			order:       application.OrderASC,
			wantIDOrder: []int64{4, 6, 10},
		},
		{
			name:        "original_url DESC",
			sortBy:      "original_url",
			order:       application.OrderDESC,
			wantIDOrder: []int64{5, 2, 8},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)

			withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
				seedDB(t, tx, "links")

				linksStore := NewLinksStore(tx)

				sort, err := application.NewSortOrder(tc.sortBy, string(tc.order))
				require.NoError(t, err)

				got, err := linksStore.GetLinks(ctx, application.GetLinksParams{
					Range: testutils.Range(t, 0, 2),
					Sort:  sort,
				})

				require.NoError(t, err)
				require.Equal(t, tc.wantIDOrder, linksIDs(got))
			})
		})
	}
}

func TestGetLinksErrors(t *testing.T) {
	validRange, err := application.NewRange(0, 5)
	require.NoError(t, err)

	validSort, err := application.NewSortOrder("id", application.OrderASC)
	require.NoError(t, err)

	t.Run("unsupported sort column", func(t *testing.T) {
		db := setupTestDB(t)

		withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
			linksStore := NewLinksStore(tx)

			sortOrder, err := application.NewSortOrder("does_not_exist", application.OrderASC)
			require.NoError(t, err)

			_, err = linksStore.GetLinks(ctx, application.GetLinksParams{
				Range: validRange,
				Sort:  sortOrder,
			})

			require.ErrorIs(t, err, application.ErrUnsupportedSortOrder)
		})
	})

	t.Run("internal store error", func(t *testing.T) {
		db := setupTestDB(t)

		withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
			linksStore := NewLinksStore(tx)

			err := tx.Rollback()
			require.NoError(t, err)

			_, err = linksStore.GetLinks(ctx, application.GetLinksParams{
				Range: validRange,
				Sort:  validSort,
			})

			require.ErrorIs(t, err, application.ErrStoreInternal)
		})
	})

	t.Run("invalid store value", func(t *testing.T) {
		db := setupTestDB(t)

		withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
			linksStore := NewLinksStore(tx)

			// NOTE: bypass store to write invalid value
			_, err := linksStore.q.CreateLink(ctx, sqlcgen.CreateLinkParams{
				OriginalUrl: "noturl",
				Shortcode:   "test",
			})
			require.NoError(t, err)

			_, err = linksStore.GetLinks(ctx, application.GetLinksParams{
				Range: validRange,
				Sort:  validSort,
			})

			require.ErrorIs(t, err, application.ErrInvalidStoreValue)
		})
	})
}

func TestGetLinksTotalCount(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		db := setupTestDB(t)

		withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
			store := NewLinksStore(tx)

			got, err := store.GetLinksTotalCount(ctx)
			require.NoError(t, err)

			require.Equal(t, int64(0), got)
		})
	})

	t.Run("with links", func(t *testing.T) {
		db := setupTestDB(t)

		withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
			store := NewLinksStore(tx)

			seedDB(t, tx, "links")

			got, err := store.GetLinksTotalCount(ctx)
			require.NoError(t, err)

			require.Equal(t, int64(10), got)
		})
	})
}

func TestUpdateLink(t *testing.T) {
	db := setupTestDB(t)

	t.Run("create and update", func(t *testing.T) {
		withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
			linksStore := NewLinksStore(tx)

			link := testutils.Link(t, "test")

			createdLink, err := linksStore.CreateLink(ctx, link)
			require.NoError(t, err)

			update := testutils.Link(t, "test-2")
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

			link := testutils.Link(t, "test")

			_, err := linksStore.UpdateLink(ctx, link)
			require.ErrorIs(t, err, application.ErrLinkNotFound)
		})
	})

	t.Run("error on shortcode change, conflict", func(t *testing.T) {
		withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
			linksStore := NewLinksStore(tx)

			link1 := testutils.Link(t, "test-1")
			link2 := testutils.Link(t, "test-2")

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

			link := testutils.Link(t, "test")

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

		link := testutils.Link(t, "test")

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

		link := testutils.Link(t, "test")

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

		_, err = linksStore.GetLinks(ctx, application.GetLinksParams{
			Range: testutils.Range(t, 0, 4),
			Sort:  testutils.SortOrder(t, "id", "ASC"),
		})
		require.ErrorIs(t, err, application.ErrInvalidStoreValue)
		require.ErrorIs(t, err, links.ErrInvlalidShortCode)
	})
}

func TestErrorInteralStore(t *testing.T) {
	db := setupTestDB(t)

	withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
		linksStore := NewLinksStore(tx)

		link := testutils.Link(t, "test")

		err := tx.Rollback()
		require.NoError(t, err)

		_, err = linksStore.CreateLink(ctx, link)
		require.ErrorIs(t, err, application.ErrStoreInternal)
	})
}
