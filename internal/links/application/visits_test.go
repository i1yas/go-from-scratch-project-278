package application

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
)

func TestGetVisits(t *testing.T) {
	cases := []struct {
		name      string
		rang      Range
		sort      SortOrder
		sortBy    string
		order     string
		setup     func(store *fakeVisitsStore)
		wantTotal int64
		wantItems []links.Visit
		err       error
	}{
		{
			name: "basic case",
			rang: createRange(t, 0, 2),
			sort: createSort(t, "id", "ASC"),
			setup: func(store *fakeVisitsStore) {
				store.
					On("GetVisitsTotalCount", mock.Anything).
					Return(int64(10), nil).
					On("GetVisits", mock.Anything, GetVisitsParams{
						Range: createRange(t, 0, 2),
						Sort:  createSort(t, "id", "ASC"),
					}).
					Return(createVisits(t, 3), nil)
			},
			wantTotal: 10,
			wantItems: createVisits(t, 3),
		},
		{
			name: "no visits",
			rang: createRange(t, 0, 2),
			sort: createSort(t, "id", "ASC"),
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
			rang: createRange(t, 0, 2),
			sort: createSort(t, "unknown", "ASC"),
			setup: func(store *fakeVisitsStore) {
				store.
					On("GetVisitsTotalCount", mock.Anything).
					Return(int64(10), nil).
					On("GetVisits", mock.Anything, GetVisitsParams{
						Range: createRange(t, 0, 2),
						Sort:  createSort(t, "unknown", "ASC"),
					}).
					Return([]links.Visit{}, ErrUnsupportedSortOrder)
			},
			err: ErrUnsupportedSortOrder,
		},
		{
			name: "reffer remaped",
			rang: createRange(t, 0, 2),
			sort: createSort(t, "reffer", "ASC"),
			setup: func(store *fakeVisitsStore) {
				store.
					On("GetVisitsTotalCount", mock.Anything).
					Return(int64(10), nil).
					On("GetVisits", mock.Anything, GetVisitsParams{
						Range: createRange(t, 0, 2),
						Sort:  createSort(t, "referer", "ASC"),
					}).
					Return(createVisits(t, 3), nil)
			},
			wantTotal: 10,
			wantItems: createVisits(t, 3),
		},
		{
			name: "store error on total count",
			rang: createRange(t, 0, 2),
			sort: createSort(t, "id", "ASC"),
			setup: func(store *fakeVisitsStore) {
				store.
					On("GetVisitsTotalCount", mock.Anything).
					Return(int64(0), ErrStoreInternal)
			},
			err: ErrStoreInternal,
		},
		{
			name: "store error on items loading",
			rang: createRange(t, 0, 2),
			sort: createSort(t, "id", "ASC"),
			setup: func(store *fakeVisitsStore) {
				store.
					On("GetVisitsTotalCount", mock.Anything).
					Return(int64(10), nil).
					On("GetVisits", mock.Anything, mock.Anything).
					Return([]links.Visit{}, ErrStoreInternal)
			},
			err: ErrStoreInternal,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := new(fakeVisitsStore)

			tc.setup(store)

			svc := NewService(&fakeLinksStore{}, store, &fakeShortcodeGen{})

			got, err := svc.GetVisits(t.Context(), GetVisitsParams{
				Range: tc.rang,
				Sort:  tc.sort,
			})

			require.ErrorIs(t, err, tc.err)
			require.Equal(t, tc.wantTotal, got.Total)
			require.Equal(t, tc.wantItems, got.Items)
		})
	}
}
