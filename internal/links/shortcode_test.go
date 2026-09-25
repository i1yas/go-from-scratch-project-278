package links

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShortCode(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		isValid bool
	}{
		{
			name:    "valid alpha shortcode",
			input:   "test",
			isValid: true,
		},
		{
			name:    "valid num shortcode",
			input:   "123",
			isValid: true,
		},
		{
			name:    "valid alphanum shortcode",
			input:   "test123",
			isValid: true,
		},
		{
			name:    "too short shortcode",
			input:   "ab",
			isValid: false,
		},
		{
			name:    "too long shortcode",
			input:   strings.Repeat("a", 33),
			isValid: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewShortCode(tc.input)

			if tc.isValid {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, ErrInvlalidShortCode)
			}
		})
	}
}
