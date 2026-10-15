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
			name:    "valid alpha shortcode with upper case",
			input:   "TEST",
			isValid: true,
		},
		{
			name:    "valid alpha shortcode with unicode",
			input:   "тест",
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
			name:    "valid with non-alphanum characters",
			input:   "-_.@",
			isValid: true,
		},
		{
			name:    "invalid with slash",
			input:   "test/test",
			isValid: false,
		},
		{
			name:    "invalid with %",
			input:   "test%test",
			isValid: false,
		},
		{
			name:    "invalid with #",
			input:   "test#test",
			isValid: false,
		},
		{
			name:    "invalid with ?",
			input:   "test?test",
			isValid: false,
		},
		{
			name:    "invalid with space",
			input:   "test test",
			isValid: false,
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
