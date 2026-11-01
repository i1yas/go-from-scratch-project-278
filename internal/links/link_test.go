package links

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLink(t *testing.T) {
	cases := []struct {
		name        string
		originalURL string
		shortcode   string
		want        Link
		err         error
	}{
		{
			name:        "valid link",
			originalURL: "http://test.com",
			shortcode:   "test",
			want: Link{
				OriginalURL: URL("http://test.com"),
				ShortCode:   ShortCode("test"),
			},
		},
		{
			name:        "error, empty url",
			originalURL: "",
			shortcode:   "test",
			err:         ErrInvlalidURL,
		},
		{
			name:        "error, empty shortcode",
			originalURL: "http://test.com",
			shortcode:   "",
			err:         ErrInvlalidShortCode,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewLink(tc.originalURL, tc.shortcode)
			require.Equal(t, tc.want, got)
			require.ErrorIs(t, err, tc.err)
		})
	}
}
