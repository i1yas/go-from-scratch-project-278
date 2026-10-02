package application

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
)

func TestGetVisitsPagination(t *testing.T) {
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
			visits := make([]links.Visit, tc.resultCount)

			for i := range tc.resultCount {
				visit, err := links.NewVisit(
					0, "0.0.0.0", "", "", 302,
				)
				require.NoError(t, err)

				visits[i] = visit
			}

			listRange, err := NewRange(tc.from, tc.to)
			require.NoError(t, err)

			params := GetVisitsParams{
				Range: listRange,
			}

			store := new(fakeVisitsStore)

			store.
				On("GetVisits", mock.Anything, params).
				Return(visits, nil).
				On("GetVisitsTotalCount", mock.Anything).
				Return(tc.total, nil)

			generator := &fakeShortcodeGen{}
			svc := NewService(
				&fakeLinksStore{},
				store,
				generator,
			)

			got, err := svc.GetVisits(t.Context(), params)

			require.NoError(t, err)
			require.Equal(t, tc.total, got.Total)
			require.Equal(t, tc.storeCalls, len(store.Calls))
			require.Equal(t, tc.resultCount, len(got.Items))

			if tc.resultCount > 0 {
				require.Equal(t, visits, got.Items)
			}
		})
	}

	// TODO: test errors
}
