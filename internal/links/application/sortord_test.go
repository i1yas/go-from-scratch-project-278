package application

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSortOrder(t *testing.T) {
	cases := []struct {
		name      string
		sortBy    string
		direction string
		want      string
		err       error
	}{
		{
			name:      "id ASC",
			sortBy:    "id",
			direction: "ASC",
			want:      "id ASC",
		},
		{
			name:      "id DESC",
			sortBy:    "id",
			direction: "DESC",
			want:      "id DESC",
		},
		{
			name:      "unknown direction",
			sortBy:    "id",
			direction: "NONE",
			err:       ErrInvalidSortOrder,
		},
		{
			name:      "empty direction",
			sortBy:    "id",
			direction: "",
			err:       ErrInvalidSortOrder,
		},
		{
			name:      "empty field",
			sortBy:    "",
			direction: "ASC",
			err:       ErrInvalidSortOrder,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewSortOrder(tc.sortBy, tc.direction)

			require.ErrorIs(t, err, tc.err)
			require.Equal(t, tc.want, got.String())
		})
	}
}
