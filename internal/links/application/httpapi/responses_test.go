package httpapi

import (
	"testing"

	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links/application/testutils"
)

func TestFormatContentRangeHeader(t *testing.T) {
	cases := []struct {
		name   string
		params contentRangeParams
		want   string
	}{
		{
			name: "first 5 items, got all 5",
			params: contentRangeParams{
				itemName:   "items",
				itemsRange: testutils.Range(t, 0, 4),
				itemsCount: 5,
				totalCount: 10,
			},
			want: "items 0-4/10",
		},
		{
			name: "first 5 items, got less",
			params: contentRangeParams{
				itemName:   "items",
				itemsRange: testutils.Range(t, 0, 4),
				itemsCount: 3,
				totalCount: 3,
			},
			want: "items 0-2/3",
		},
		{
			name: "from 5th to 8th, got all 4",
			params: contentRangeParams{
				itemName:   "items",
				itemsRange: testutils.Range(t, 5, 8),
				itemsCount: 4,
				totalCount: 10,
			},
			want: "items 5-8/10",
		},
		{
			name: "first 5 items, got none of 0",
			params: contentRangeParams{
				itemName:   "items",
				itemsRange: testutils.Range(t, 0, 4),
				itemsCount: 0,
				totalCount: 0,
			},
			want: "items */0",
		},
		{
			name: "from 800 to 810 items, got none of 10",
			params: contentRangeParams{
				itemName:   "items",
				itemsRange: testutils.Range(t, 800, 810),
				itemsCount: 0,
				totalCount: 10,
			},
			want: "items */10",
		},
	}

	for _, tc := range cases {
		got := formatContentRangeHeader(tc.params)
		require.Equal(t, tc.want, got)
	}
}
