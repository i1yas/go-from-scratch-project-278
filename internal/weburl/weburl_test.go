package weburl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
		err   error
	}{
		{
			name:  "basic case with http",
			input: "http://test.com",
			want:  "http://test.com",
		},
		{
			name:  "basic case with https",
			input: "https://test.com",
			want:  "https://test.com",
		},
		{
			name:  "host normamlized to lower case",
			input: "http://TEST.com",
			want:  "http://test.com",
		},
		{
			name:  "url without tld is valid",
			input: "http://test",
			want:  "http://test",
		},
		{
			name:  "url with port in host is valid",
			input: "http://localhost:8080",
			want:  "http://localhost:8080",
		},
		{
			name:  "url without host is invalid",
			input: "http://",
			err:   ErrMissingHost,
		},
		{
			name:  "failed to parse url",
			input: "://test.com",
			err:   ErrFailedToParse,
		},
		{
			name:  "unsupported protocol",
			input: "ws://test.com",
			err:   ErrUnsupportedProtocol,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Parse(tc.input)
			require.ErrorIs(t, err, tc.err)

			if tc.err == nil {
				require.Equal(t, tc.want, got.String())
			}
		})
	}
}
