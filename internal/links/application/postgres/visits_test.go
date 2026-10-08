package postgres

import (
	"context"
	"database/sql"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
	"hexleturlshort/internal/links/application/testutils"
)

func TestGetVisitsRange(t *testing.T) {
	cases := []struct {
		name      string
		rang      application.Range
		wantItems []int64
	}{
		{
			name:      "start, from 0 to 10",
			rang:      testutils.Range(t, 0, 9),
			wantItems: []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:      "middle, from 4 to 14",
			rang:      testutils.Range(t, 4, 13),
			wantItems: []int64{5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
		},
		{
			name:      "end, from 10 to 20",
			rang:      testutils.Range(t, 10, 19),
			wantItems: []int64{11, 12, 13, 14, 15},
		},
	}

	db := setupTestDB(t)

	withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
		seedDB(t, tx, "links")
		seedDB(t, tx, "visits")

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				store := NewVisitsStore(tx)

				sort, err := application.NewSortOrder("id", "ASC")
				require.NoError(t, err)

				got, err := store.GetVisits(ctx, application.GetVisitsParams{
					Range: tc.rang,
					Sort:  sort,
				})
				require.NoError(t, err)

				require.Equal(t, len(tc.wantItems), len(got))

				for i := range tc.wantItems {
					require.Equal(t, tc.wantItems[i], got[i].ID)
				}
			})
		}
	})
}

func TestGetVisitsSort(t *testing.T) {
	cases := []struct {
		name      string
		sortBy    string
		order     string
		wantItems []int64
	}{
		{
			name:      "sort by id ASC",
			sortBy:    "id",
			order:     "ASC",
			wantItems: []int64{1, 2, 3, 4},
		},
		{
			name:      "sort by id DESC",
			sortBy:    "id",
			order:     "DESC",
			wantItems: []int64{15, 14, 13, 12},
		},
		{
			name:      "sort by link_id ASC",
			sortBy:    "link_id",
			order:     "ASC",
			wantItems: []int64{1, 10, 4, 15},
		},
		{
			name:      "sort by link_id DESC",
			sortBy:    "link_id",
			order:     "DESC",
			wantItems: []int64{14, 3, 6, 8},
		},
		{
			name:      "sort by IP ASC",
			sortBy:    "ip",
			order:     "ASC",
			wantItems: []int64{7, 1, 11, 9},
		},
		{
			name:      "sort by IP DESC",
			sortBy:    "ip",
			order:     "DESC",
			wantItems: []int64{8, 10, 4, 14},
		},
		{
			name:      "sort by referer ASC",
			sortBy:    "referer",
			order:     "ASC",
			wantItems: []int64{7, 13, 12, 8},
		},
		{
			name:      "sort by referer DESC",
			sortBy:    "referer",
			order:     "DESC",
			wantItems: []int64{11, 10, 4, 6},
		},
		{
			name:      "sort by user_agent ASC",
			sortBy:    "user_agent",
			order:     "ASC",
			wantItems: []int64{4, 7, 12, 14},
		},
		{
			name:      "sort by user_agent DESC",
			sortBy:    "user_agent",
			order:     "DESC",
			wantItems: []int64{1, 2, 3, 5},
		},
		{
			name:      "sort by status ASC",
			sortBy:    "status",
			order:     "ASC",
			wantItems: []int64{1, 2, 3, 4},
		},
		{
			name:      "sort by status DESC",
			sortBy:    "status",
			order:     "DESC",
			wantItems: []int64{5, 6, 9, 13},
		},
		{
			name:      "sort by created_at ASC",
			sortBy:    "created_at",
			order:     "ASC",
			wantItems: []int64{6, 12, 7, 2},
		},
		{
			name:      "sort by created_at DESC",
			sortBy:    "created_at",
			order:     "DESC",
			wantItems: []int64{14, 9, 4, 10},
		},
	}

	db := setupTestDB(t)

	withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
		seedDB(t, tx, "links")
		seedDB(t, tx, "visits")

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				store := NewVisitsStore(tx)

				rang, err := application.NewRange(0, 3)
				require.NoError(t, err)

				sort, err := application.NewSortOrder(tc.sortBy, tc.order)
				require.NoError(t, err)

				got, err := store.GetVisits(ctx, application.GetVisitsParams{
					Range: rang,
					Sort:  sort,
				})
				require.NoError(t, err)

				// if tc.sortBy == "link_id" {
				// 	fmt.Printf("sort=%v, got=%v\n", sort, got)
				// 	t.Fatal("DEBUG")
				// }

				require.Equal(t, len(tc.wantItems), len(got))
				require.Equal(t, tc.wantItems, visitsIDs(got))
			})
		}
	})
}

func TestGetVisitsFormatting(t *testing.T) {
	db := setupTestDB(t)

	withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
		seedDB(t, tx, "links")
		seedDB(t, tx, "visits")

		store := NewVisitsStore(tx)
		rang, err := application.NewRange(0, 6)
		require.NoError(t, err)

		sort, err := application.NewSortOrder("id", "ASC")
		require.NoError(t, err)

		visits, err := store.GetVisits(ctx, application.GetVisitsParams{
			Range: rang,
			Sort:  sort,
		})
		require.NoError(t, err)

		wantVisit1 := links.Visit{
			ID:        1,
			LinkID:    1,
			IP:        netip.AddrFrom4([4]byte{36, 8, 244, 30}),
			Referer:   "https://example.org/nested/page",
			UserAgent: "very-very-long-user-agent",
			Status:    302,
			CreatedAt: timeFromISO(t, "2026-09-19T07:56:07.602Z"),
		}

		wantVisit7 := links.Visit{
			ID:        7,
			LinkID:    5,
			IP:        netip.AddrFrom4([4]byte{0, 49, 240, 5}),
			Referer:   "",
			UserAgent: "",
			Status:    302,
			CreatedAt: timeFromISO(t, "2026-08-17T06:42:16.056Z"),
		}

		requireEqualVisit(t, wantVisit1, visits[0])
		requireEqualVisit(t, wantVisit7, visits[6])
	})
}

func TestGetVisitsErrors(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		t.Run("internal error", func(t *testing.T) {
			db := setupTestDB(t)

			withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
				store := NewVisitsStore(tx)

				err := tx.Rollback()
				require.NoError(t, err)

				rang, err := application.NewRange(0, 4)
				require.NoError(t, err)

				sort, err := application.NewSortOrder("id", "ASC")
				require.NoError(t, err)

				_, err = store.GetVisits(ctx, application.GetVisitsParams{
					Range: rang,
					Sort:  sort,
				})
				require.ErrorIs(t, err, application.ErrStoreInternal)
			})
		})

		t.Run("unsupported sort", func(t *testing.T) {
			db := setupTestDB(t)

			withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
				store := NewVisitsStore(tx)

				rang, err := application.NewRange(0, 4)
				require.NoError(t, err)

				sort, err := application.NewSortOrder("somefield", "ASC")
				require.NoError(t, err)

				_, err = store.GetVisits(ctx, application.GetVisitsParams{
					Range: rang,
					Sort:  sort,
				})
				require.ErrorIs(t, err, application.ErrUnsupportedSortOrder)
			})
		})
	})
}

func TestGetVisitsTotalCount(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		db := setupTestDB(t)

		withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
			store := NewVisitsStore(tx)

			got, err := store.GetVisitsTotalCount(ctx)
			require.NoError(t, err)

			require.Equal(t, int64(0), got)
		})
	})

	t.Run("with visits", func(t *testing.T) {
		db := setupTestDB(t)

		withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
			store := NewVisitsStore(tx)

			seedDB(t, tx, "links")
			seedDB(t, tx, "visits")

			got, err := store.GetVisitsTotalCount(ctx)
			require.NoError(t, err)

			require.Equal(t, int64(15), got)
		})
	})
}

func TestCreateVisit(t *testing.T) {
	visitWithOptionals, err := links.NewVisit(
		1,
		"1.2.3.4",
		"http://test.com",
		"user-agent",
		302,
	)
	require.NoError(t, err)

	visitWithoutOptionals, err := links.NewVisit(
		1,
		"1.2.3.4",
		"http://test.com",
		"user-agent",
		302,
	)
	require.NoError(t, err)

	cases := []struct {
		name  string
		input links.Visit
	}{
		{
			name:  "with optional fields",
			input: visitWithOptionals,
		},
		{
			name:  "without optional fields",
			input: visitWithoutOptionals,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)

			withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
				store := NewVisitsStore(tx)

				seedDB(t, tx, "links")

				err = store.CreateVisit(ctx, tc.input)
				require.NoError(t, err)

				rang, err := application.NewRange(0, 1)
				require.NoError(t, err)

				sort, err := application.NewSortOrder("id", "ASC")
				require.NoError(t, err)

				visits, err := store.GetVisits(ctx, application.GetVisitsParams{
					Range: rang,
					Sort:  sort,
				})
				require.NoError(t, err)

				got := visits[0]

				require.Equal(t, tc.input.LinkID, got.LinkID)
				require.Equal(t, tc.input.IP, got.IP)
				require.Equal(t, tc.input.Referer, got.Referer)
				require.Equal(t, tc.input.UserAgent, got.UserAgent)
				require.Equal(t, tc.input.Status, got.Status)
			})
		})
	}

	t.Run("link does not exist", func(t *testing.T) {
		db := setupTestDB(t)

		withTx(t, db, func(ctx context.Context, tx *sql.Tx) {
			store := NewVisitsStore(tx)

			err := store.CreateVisit(ctx, visitWithOptionals)
			require.ErrorIs(t, err, application.ErrStoreInternal)
		})
	})
}
