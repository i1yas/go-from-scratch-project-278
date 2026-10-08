package application_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
	"hexleturlshort/internal/links/application/testutils"
)

func TestGetVisits(t *testing.T) {
	cases := []struct {
		name      string
		rang      application.Range
		sort      application.SortOrder
		sortBy    string
		order     string
		setup     func(store *fakeVisitsStore)
		wantTotal int64
		wantItems []links.Visit
		err       error
	}{
		{
			name: "basic case",
			rang: testutils.Range(t, 0, 2),
			sort: testutils.SortOrder(t, "id", "ASC"),
			setup: func(store *fakeVisitsStore) {
				store.
					On("GetVisitsTotalCount", mock.Anything).
					Return(int64(10), nil).
					On("GetVisits", mock.Anything, application.GetVisitsParams{
						Range: testutils.Range(t, 0, 2),
						Sort:  testutils.SortOrder(t, "id", "ASC"),
					}).
					Return(testutils.Visits(t, 3), nil)
			},
			wantTotal: 10,
			wantItems: testutils.Visits(t, 3),
		},
		{
			name: "no visits",
			rang: testutils.Range(t, 0, 2),
			sort: testutils.SortOrder(t, "id", "ASC"),
			setup: func(store *fakeVisitsStore) {
				store.
					On("GetVisitsTotalCount", mock.Anything).
					Return(int64(0), nil)
			},
			wantTotal: 0,
			wantItems: nil,
		},
		{
			name: "unsupported sorting",
			rang: testutils.Range(t, 0, 2),
			sort: testutils.SortOrder(t, "unknown", "ASC"),
			setup: func(store *fakeVisitsStore) {
				store.
					On("GetVisitsTotalCount", mock.Anything).
					Return(int64(10), nil).
					On("GetVisits", mock.Anything, application.GetVisitsParams{
						Range: testutils.Range(t, 0, 2),
						Sort:  testutils.SortOrder(t, "unknown", "ASC"),
					}).
					Return([]links.Visit{}, application.ErrUnsupportedSortOrder)
			},
			err: application.ErrUnsupportedSortOrder,
		},
		{
			name: "reffer remaped",
			rang: testutils.Range(t, 0, 2),
			sort: testutils.SortOrder(t, "reffer", "ASC"),
			setup: func(store *fakeVisitsStore) {
				store.
					On("GetVisitsTotalCount", mock.Anything).
					Return(int64(10), nil).
					On("GetVisits", mock.Anything, application.GetVisitsParams{
						Range: testutils.Range(t, 0, 2),
						Sort:  testutils.SortOrder(t, "referer", "ASC"),
					}).
					Return(testutils.Visits(t, 3), nil)
			},
			wantTotal: 10,
			wantItems: testutils.Visits(t, 3),
		},
		{
			name: "store error on total count",
			rang: testutils.Range(t, 0, 2),
			sort: testutils.SortOrder(t, "id", "ASC"),
			setup: func(store *fakeVisitsStore) {
				store.
					On("GetVisitsTotalCount", mock.Anything).
					Return(int64(0), application.ErrStoreInternal)
			},
			err: application.ErrStoreInternal,
		},
		{
			name: "store error on items loading",
			rang: testutils.Range(t, 0, 2),
			sort: testutils.SortOrder(t, "id", "ASC"),
			setup: func(store *fakeVisitsStore) {
				store.
					On("GetVisitsTotalCount", mock.Anything).
					Return(int64(10), nil).
					On("GetVisits", mock.Anything, mock.Anything).
					Return([]links.Visit{}, application.ErrStoreInternal)
			},
			err: application.ErrStoreInternal,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := new(fakeVisitsStore)

			tc.setup(store)

			svc := application.NewService(&fakeLinksStore{}, store, &fakeShortcodeGen{})

			got, err := svc.GetVisits(t.Context(), application.GetVisitsParams{
				Range: tc.rang,
				Sort:  tc.sort,
			})

			require.ErrorIs(t, err, tc.err)
			require.Equal(t, tc.wantTotal, got.Total)
			require.Equal(t, tc.wantItems, got.Items)
		})
	}
}
