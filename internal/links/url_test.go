package links

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestURL(t *testing.T) {
	maxLenURL := "http://test/" + strings.Repeat("x", 188)
	maxLenWithCyrillicHost := "http://пример.рф/" + strings.Repeat("x", 171)

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
			name:    "valid url unicode",
			input:   "https://пример.рф",
			isValid: true,
		},
		{
			name:    "max length url is valid",
			input:   maxLenURL,
			isValid: true,
		},
		{
			name:    "max length url with cyrillic host is valid",
			input:   maxLenWithCyrillicHost,
			isValid: true,
		},
		{
			name:    "empty url is invalid",
			input:   "",
			isValid: false,
		},
		{
			name:    "url is too long",
			input:   maxLenURL + "x",
			isValid: false,
		},
		{
			name:    "url with cyrillic host is too long",
			input:   maxLenWithCyrillicHost + "x",
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
