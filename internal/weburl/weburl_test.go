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
			name:  "basic case with subdomain",
			input: "https://subtest.test.com",
			want:  "https://subtest.test.com",
		},
		{
			name:  "host normamlized to lower case",
			input: "http://TEST.COM",
			want:  "http://test.com",
		},
		{
			name:  "cyrillic host converted to punycode",
			input: "https://пример.рф/page",
			want:  "https://xn--e1afmkfd.xn--p1ai/page",
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
			name:  "url with ipv4 is invalid",
			input: "http://1.2.3.4",
			err:   ErrFailedToParse,
		},
		{
			name:  "url with ipv6 is invalid",
			input: "http://2001:db8::1",
			err:   ErrFailedToParse,
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
