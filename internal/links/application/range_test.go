package application

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRange(t *testing.T) {
	cases := []struct {
		name string
		from int
		to   int
		want Range
		err  error
	}{
		{
			name: "from 0 to 5",
			from: 0,
			to:   5,
			want: Range{0, 5},
		},
		{
			name: "from 0 to 0",
			from: 0,
			to:   0,
			want: Range{0, 0},
		},
		{
			name: "from 2 to 5",
			from: 2,
			to:   5,
			want: Range{2, 5},
		},
		{
			name: "error from -5 to 5",
			from: -5,
			to:   5,
			err:  ErrInvalidRange,
		},
		{
			name: "error from 0 to -5",
			from: 0,
			to:   -5,
			err:  ErrInvalidRange,
		},
		{
			name: "error from 5 to 3",
			from: 5,
			to:   3,
			err:  ErrInvalidRange,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewRange(int32(tc.from), int32(tc.to))

			require.ErrorIs(t, err, tc.err)
			require.Equal(t, tc.want, got)
		})
	}
}
