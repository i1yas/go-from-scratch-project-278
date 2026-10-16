package httpapi

import (
	"testing"

	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links/application"
)

func TestParseRange(t *testing.T) {
	fallback := application.Range{From: 0, To: 9}

	cases := []struct {
		name        string
		input       string
		maxPageSize int32
		want        application.Range
		err         error
	}{
		{
			name:        "valid range",
			input:       "[0,4]",
			maxPageSize: 100,
			want:        application.Range{From: 0, To: 4},
		},
		{
			name:        "valid range with spaces",
			input:       "[ 0 , 4 ]",
			maxPageSize: 100,
			want:        application.Range{From: 0, To: 4},
		},
		{
			name:        "fallback on empty input",
			input:       "",
			maxPageSize: 100,
			want:        fallback,
		},
		{
			name:        "skip extra elements",
			input:       "[0,4, -1]",
			maxPageSize: 100,
			want:        application.Range{From: 0, To: 4},
		},
		{
			name:        "lower bound less than 0",
			input:       "[-1,4]",
			maxPageSize: 100,
			err:         application.ErrInvalidRange,
		},
		{
			name:        "upper bound less lower bound",
			input:       "[3,1]",
			maxPageSize: 100,
			err:         application.ErrInvalidRange,
		},
		{
			name:        "uppper bound missing",
			input:       "[0]",
			maxPageSize: 100,
			err:         ErrInvalidQuery,
		},
		{
			name:        "invalid format",
			input:       "1,2",
			maxPageSize: 100,
			err:         ErrInvalidQuery,
		},
		{
			name:        "range is almost too big",
			input:       "[0,99]",
			maxPageSize: 100,
			want:        application.Range{From: 0, To: 99},
		},
		{
			name:        "range is too big",
			input:       "[0,100]",
			maxPageSize: 100,
			err:         ErrRangeIsTooBig,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseRange(tc.input, 100, fallback)

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
			name:  "missing element",
			input: `["name"]`,
			err:   ErrInvalidQuery,
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
