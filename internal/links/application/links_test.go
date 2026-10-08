package application_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
	"hexleturlshort/internal/links/application/testutils"
)

func TestGetLinks(t *testing.T) {
	cases := []struct {
		name      string
		rang      application.Range
		sort      application.SortOrder
		sortBy    string
		order     string
		setup     func(store *fakeLinksStore)
		wantTotal int64
		wantItems []links.Link
		err       error
	}{
		{
			name: "basic case",
			rang: testutils.Range(t, 0, 2),
			sort: testutils.SortOrder(t, "id", "ASC"),
			setup: func(store *fakeLinksStore) {
				store.
					On("GetLinksTotalCount", mock.Anything).
					Return(int64(10), nil).
					On("GetLinks", mock.Anything, application.GetLinksParams{
						Range: testutils.Range(t, 0, 2),
						Sort:  testutils.SortOrder(t, "id", "ASC"),
					}).
					Return(testutils.Links(t, 3), nil)
			},
			wantTotal: 10,
			wantItems: testutils.Links(t, 3),
		},
		{
			name: "no links",
			rang: testutils.Range(t, 0, 2),
			sort: testutils.SortOrder(t, "id", "ASC"),
			setup: func(store *fakeLinksStore) {
				store.
					On("GetLinksTotalCount", mock.Anything).
					Return(int64(0), nil)
			},
			wantTotal: 0,
			wantItems: nil,
		},
		{
			name: "unsupported sorting",
			rang: testutils.Range(t, 0, 2),
			sort: testutils.SortOrder(t, "unknown", "ASC"),
			setup: func(store *fakeLinksStore) {
				store.
					On("GetLinksTotalCount", mock.Anything).
					Return(int64(10), nil).
					On("GetLinks", mock.Anything, application.GetLinksParams{
						Range: testutils.Range(t, 0, 2),
						Sort:  testutils.SortOrder(t, "unknown", "ASC"),
					}).
					Return([]links.Link{}, application.ErrUnsupportedSortOrder)
			},
			err: application.ErrUnsupportedSortOrder,
		},
		{
			name: "short_name remapped to shortcode",
			rang: testutils.Range(t, 0, 2),
			sort: testutils.SortOrder(t, "short_name", "ASC"),
			setup: func(store *fakeLinksStore) {
				store.
					On("GetLinksTotalCount", mock.Anything).
					Return(int64(10), nil).
					On("GetLinks", mock.Anything, application.GetLinksParams{
						Range: testutils.Range(t, 0, 2),
						Sort:  testutils.SortOrder(t, "shortcode", "ASC"),
					}).
					Return(testutils.Links(t, 3), nil)
			},
			wantTotal: 10,
			wantItems: testutils.Links(t, 3),
		},
		{
			name: "short_url remapped to shortcode",
			rang: testutils.Range(t, 0, 2),
			sort: testutils.SortOrder(t, "short_url", "ASC"),
			setup: func(store *fakeLinksStore) {
				store.
					On("GetLinksTotalCount", mock.Anything).
					Return(int64(10), nil).
					On("GetLinks", mock.Anything, application.GetLinksParams{
						Range: testutils.Range(t, 0, 2),
						Sort:  testutils.SortOrder(t, "shortcode", "ASC"),
					}).
					Return(testutils.Links(t, 3), nil)
			},
			wantTotal: 10,
			wantItems: testutils.Links(t, 3),
		},
		{
			name: "store error on total count",
			rang: testutils.Range(t, 0, 2),
			sort: testutils.SortOrder(t, "id", "ASC"),
			setup: func(store *fakeLinksStore) {
				store.
					On("GetLinksTotalCount", mock.Anything).
					Return(int64(0), application.ErrStoreInternal)
			},
			err: application.ErrStoreInternal,
		},
		{
			name: "store error on items loading",
			rang: testutils.Range(t, 0, 2),
			sort: testutils.SortOrder(t, "id", "ASC"),
			setup: func(store *fakeLinksStore) {
				store.
					On("GetLinksTotalCount", mock.Anything).
					Return(int64(10), nil).
					On("GetLinks", mock.Anything, mock.Anything).
					Return([]links.Link{}, application.ErrStoreInternal)
			},
			err: application.ErrStoreInternal,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := new(fakeLinksStore)

			tc.setup(store)

			svc := application.NewService(store, &fakeVisitsStore{}, &fakeShortcodeGen{})

			got, err := svc.GetLinks(t.Context(), application.GetLinksParams{
				Range: tc.rang,
				Sort:  tc.sort,
			})

			require.ErrorIs(t, err, tc.err)
			require.Equal(t, tc.wantTotal, got.Total)
			require.Equal(t, tc.wantItems, got.Items)

			store.AssertExpectations(t)
		})
	}
}

func TestGetLinkByID(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		store := new(fakeLinksStore)

		want := testutils.Link(t, "test")
		want.ID = 101

		store.
			On("GetLinkByID", mock.Anything, want.ID).
			Return(want, nil)

		generator := &fakeShortcodeGen{}
		svc := application.NewService(
			store,
			&fakeVisitsStore{},
			generator,
		)

		got, err := svc.GetLinkByID(t.Context(), want.ID)

		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("not found", func(t *testing.T) {
		store := new(fakeLinksStore)

		store.
			On("GetLinkByID", mock.Anything, mock.Anything).
			Return(links.Link{}, application.ErrLinkNotFound)

		generator := &fakeShortcodeGen{}
		svc := application.NewService(
			store,
			&fakeVisitsStore{},
			generator,
		)

		_, err := svc.GetLinkByID(t.Context(), 101)

		require.ErrorIs(t, err, application.ErrLinkNotFound)
	})
}
