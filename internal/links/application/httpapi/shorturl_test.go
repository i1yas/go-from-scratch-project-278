package httpapi

import (
	"testing"

	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
)

func TestShortURLGenerator(t *testing.T) {
	cases := []struct {
		name     string
		baseURL  string
		shortURL string
		err      error
	}{
		{
			name:     "valid base url",
			baseURL:  "http://test.com",
			shortURL: "http://test.com/r/test-code",
		},
		{
			name:     "valid base url with uppercase host",
			baseURL:  "http://TEST.com",
			shortURL: "http://test.com/r/test-code",
		},
		{
			name:    "invalid base url with unsupported protocol",
			baseURL: "ws://test.com",
			err:     ErrInvalidBaseURL,
		},
		{
			name:    "invalid base url without protocol",
			baseURL: "invalid",
			err:     ErrInvalidBaseURL,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g, err := NewShortURLGenerator(tc.baseURL)
			require.ErrorIs(t, err, tc.err)

			if tc.err == nil {
				shortcode, err := links.NewShortCode("test-code")
				require.NoError(t, err)

				shortURL, err := g.Generate(shortcode)
				require.NoError(t, err)
				require.Equal(t, tc.shortURL, shortURL)
			}
		})
	}
}
