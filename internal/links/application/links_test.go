package application

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
)

func TestGetLinks(t *testing.T) {
	cases := []struct {
		name        string
		from        int32
		to          int32
		total       int64
		resultCount int
		storeCalls  int
	}{
		{
			name:        "from 0 to 5, total 1",
			from:        0,
			to:          5,
			total:       1,
			resultCount: 1,
			storeCalls:  2,
		},
		{
			name:        "from 0 to 5, total 3",
			from:        0,
			to:          5,
			total:       3,
			resultCount: 3,
			storeCalls:  2,
		},
		{
			name:        "from 0 to 5, total 6",
			from:        0,
			to:          5,
			total:       6,
			resultCount: 5,
			storeCalls:  2,
		},
		{
			name:        "from 0 to 5, total 0",
			from:        0,
			to:          5,
			total:       0,
			resultCount: 0,
			storeCalls:  1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			linkItems := make([]links.Link, tc.resultCount)

			for i := range tc.resultCount {
				code := fmt.Sprintf("link-%d", i)
				linkItems[i] = createValidLink(t, code)
			}

			linksRange, err := NewRange(tc.from, tc.to)
			require.NoError(t, err)

			params := GetLinksParams{
				Range: linksRange,
			}

			store := new(fakeLinksStore)

			store.
				On("GetLinks", mock.Anything, params).
				Return(linkItems, nil).
				On("GetLinksTotalCount", mock.Anything).
				Return(tc.total, nil)

			generator := &fakeShortcodeGen{}
			svc := NewService(
				store,
				&fakeVisitsStore{},
				generator,
			)

			got, err := svc.GetLinks(t.Context(), params)

			require.NoError(t, err)
			require.Equal(t, tc.total, got.Total)
			require.Equal(t, tc.storeCalls, len(store.Calls))
			require.Equal(t, tc.resultCount, len(got.Items))

			if tc.resultCount > 0 {
				require.Equal(t, linkItems, got.Items)
			}
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
