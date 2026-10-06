package httpapi

import (
	"testing"

	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links/application"
)

func TestParseRange(t *testing.T) {
	fallback := application.Range{From: 0, To: 9}

	cases := []struct {
		name  string
		input string
		want  application.Range
		err   error
	}{
		{
			name:  "valid range",
			input: "[0,4]",
			want:  application.Range{From: 0, To: 4},
		},
		{
			name:  "valid range with spaces",
			input: "[ 0 , 4 ]",
			want:  application.Range{From: 0, To: 4},
		},
		{
			name:  "fallback on empty input",
			input: "",
			want:  fallback,
		},
		{
			name:  "skip extra elements",
			input: "[0,4, -1]",
			want:  application.Range{From: 0, To: 4},
		},
		{
			name:  "lower bound less than 0",
			input: "[-1,4]",
			err:   application.ErrInvalidRange,
		},
		{
			name:  "upper bound less lower bound",
			input: "[3,1]",
			err:   application.ErrInvalidRange,
		},
		{
			name:  "invalid format",
			input: "1,2",
			err:   ErrInvalidQuery,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseRange(tc.input, fallback)

			require.ErrorIs(t, err, tc.err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestParseSort(t *testing.T) {
	fallback := application.SortOrder{
		SortBy: "id",
		Order:  application.OrderASC,
	}

	cases := []struct {
		name  string
		input string
		want  application.SortOrder
		err   error
	}{
		{
			name:  "valid sort ASC",
			input: `["name","ASC"]`,
			want:  application.SortOrder{SortBy: "name", Order: application.OrderASC},
		},
		{
			name:  "valid sort DESC",
			input: `["name","DESC"]`,
			want:  application.SortOrder{SortBy: "name", Order: application.OrderDESC},
		},
		{
			name:  "valid sort DESC with spaces",
			input: `[ "name" , "DESC" ]`,
			want:  application.SortOrder{SortBy: "name", Order: application.OrderDESC},
		},
		{
			name:  "fallback on empty input",
			input: "",
			want:  fallback,
		},
		{
			name:  "skip extra elements",
			input: `["name","DESC","extra"]`,
			want:  application.SortOrder{SortBy: "name", Order: application.OrderDESC},
		},
		{
			name:  "invalid order",
			input: `["name","UNKNOWN"]`,
			err:   application.ErrInvalidSortOrder,
		},
		{
			name:  "invalid format",
			input: "id,ASC",
			err:   ErrInvalidQuery,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseSort(tc.input, fallback)

			require.ErrorIs(t, err, tc.err)
			require.Equal(t, tc.want, got)
		})
	}
}
