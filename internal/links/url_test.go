package links

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestURL(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		isValid bool
	}{
		{
			name:    "valid url with http protocol",
			input:   "http://test.com",
			isValid: true,
		},
		{
			name:    "valid url with https protocol",
			input:   "https://test.com",
			isValid: true,
		},
		{
			name:    "empty url is invalid",
			input:   "",
			isValid: false,
		},
		{
			name:    "unsupported protocol",
			input:   "ws://test.com/test",
			isValid: false,
		},
		{
			name:    "missing host",
			input:   "http://",
			isValid: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewURL(tc.input)

			if tc.isValid {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, ErrInvlalidURL)
			}
		})
	}
}
