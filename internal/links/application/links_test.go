package application

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
)

func TestGetLinks(t *testing.T) {
	cases := []struct {
		name      string
		rang      Range
		sort      SortOrder
		sortBy    string
		order     string
		setup     func(store *fakeLinksStore)
		wantTotal int64
		wantItems []links.Link
		err       error
	}{
		{
			name: "basic case",
			rang: createRange(t, 0, 2),
			sort: createSort(t, "id", "ASC"),
			setup: func(store *fakeLinksStore) {
				store.
					On("GetLinksTotalCount", mock.Anything).
					Return(int64(10), nil).
					On("GetLinks", mock.Anything, GetLinksParams{
						Range: createRange(t, 0, 2),
						Sort:  createSort(t, "id", "ASC"),
					}).
					Return(createLinks(t, 3), nil)
			},
			wantTotal: 10,
			wantItems: createLinks(t, 3),
		},
		{
			name: "no links",
			rang: createRange(t, 0, 2),
			sort: createSort(t, "id", "ASC"),
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
			rang: createRange(t, 0, 2),
			sort: createSort(t, "unknown", "ASC"),
			setup: func(store *fakeLinksStore) {
				store.
					On("GetLinksTotalCount", mock.Anything).
					Return(int64(10), nil).
					On("GetLinks", mock.Anything, GetLinksParams{
						Range: createRange(t, 0, 2),
						Sort:  createSort(t, "unknown", "ASC"),
					}).
					Return([]links.Link{}, ErrUnsupportedSortOrder)
			},
			err: ErrUnsupportedSortOrder,
		},
		{
			name: "short_name remapped to shortcode",
			rang: createRange(t, 0, 2),
			sort: createSort(t, "short_name", "ASC"),
			setup: func(store *fakeLinksStore) {
				store.
					On("GetLinksTotalCount", mock.Anything).
					Return(int64(10), nil).
					On("GetLinks", mock.Anything, GetLinksParams{
						Range: createRange(t, 0, 2),
						Sort:  createSort(t, "shortcode", "ASC"),
					}).
					Return(createLinks(t, 3), nil)
			},
			wantTotal: 10,
			wantItems: createLinks(t, 3),
		},
		{
			name: "short_url remapped to shortcode",
			rang: createRange(t, 0, 2),
			sort: createSort(t, "short_url", "ASC"),
			setup: func(store *fakeLinksStore) {
				store.
					On("GetLinksTotalCount", mock.Anything).
					Return(int64(10), nil).
					On("GetLinks", mock.Anything, GetLinksParams{
						Range: createRange(t, 0, 2),
						Sort:  createSort(t, "shortcode", "ASC"),
					}).
					Return(createLinks(t, 3), nil)
			},
			wantTotal: 10,
			wantItems: createLinks(t, 3),
		},
		{
			name: "store error on total count",
			rang: createRange(t, 0, 2),
			sort: createSort(t, "id", "ASC"),
			setup: func(store *fakeLinksStore) {
				store.
					On("GetLinksTotalCount", mock.Anything).
					Return(int64(0), ErrStoreInternal)
			},
			err: ErrStoreInternal,
		},
		{
			name: "store error on items loading",
			rang: createRange(t, 0, 2),
			sort: createSort(t, "id", "ASC"),
			setup: func(store *fakeLinksStore) {
				store.
					On("GetLinksTotalCount", mock.Anything).
					Return(int64(10), nil).
					On("GetLinks", mock.Anything, mock.Anything).
					Return([]links.Link{}, ErrStoreInternal)
			},
			err: ErrStoreInternal,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := new(fakeLinksStore)

			tc.setup(store)

			svc := NewService(store, &fakeVisitsStore{}, &fakeShortcodeGen{})

			got, err := svc.GetLinks(t.Context(), GetLinksParams{
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

		want := createValidLink(t, "test")
		want.ID = 101

		store.
			On("GetLinkByID", mock.Anything, want.ID).
			Return(want, nil)

		generator := &fakeShortcodeGen{}
		svc := NewService(
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
			Return(links.Link{}, ErrLinkNotFound)

		generator := &fakeShortcodeGen{}
		svc := NewService(
			store,
			&fakeVisitsStore{},
			generator,
		)

		_, err := svc.GetLinkByID(t.Context(), 101)

		require.ErrorIs(t, err, ErrLinkNotFound)
	})
}
