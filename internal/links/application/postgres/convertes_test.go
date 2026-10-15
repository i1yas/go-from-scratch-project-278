package postgres

import (
	"testing"

	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links/application"
	"hexleturlshort/internal/links/application/testutils"
)

func TestConvertRangeToLimitOffset(t *testing.T) {
	cases := []struct {
		name       string
		input      application.Range
		wantLimit  int32
		wantOffset int32
	}{
		{
			name:       "from 0 to 4",
			input:      testutils.Range(t, 0, 4),
			wantLimit:  5,
			wantOffset: 0,
		},
		{
			name:       "from 2 to 5",
			input:      testutils.Range(t, 2, 5),
			wantLimit:  4,
			wantOffset: 2,
		},
		{
			name:       "from 0 to 0",
			input:      testutils.Range(t, 0, 0),
			wantLimit:  1,
			wantOffset: 0,
		},
		{
			name:       "from 5 to 5",
			input:      testutils.Range(t, 5, 5),
			wantLimit:  1,
			wantOffset: 5,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotLimit, gotOffset := convertRangeToLimitOffset(tc.input)

			require.Equal(t, tc.wantLimit, gotLimit)
			require.Equal(t, tc.wantOffset, gotOffset)
		})
	}
}
