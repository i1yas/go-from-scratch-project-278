package links

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLink(t *testing.T) {
	cases := []struct {
		name        string
		originalURL string
		shortcode   string
		want        Link
		fieldsErrs  map[string]string
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
			fieldsErrs: map[string]string{
				"original_url": ErrInvlalidURL.Error(),
			},
		},
		{
			name:        "error, empty shortcode",
			originalURL: "http://test.com",
			shortcode:   "",
			fieldsErrs: map[string]string{
				"shortcode": ErrInvlalidShortCode.Error(),
			},
		},
		{
			name:        "error, both original url and shortcode are empty",
			originalURL: "",
			shortcode:   "",
			fieldsErrs: map[string]string{
				"original_url": ErrInvlalidURL.Error(),
				"shortcode":    ErrInvlalidShortCode.Error(),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewLink(tc.originalURL, tc.shortcode)

			if tc.fieldsErrs == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				linkErr, ok := errors.AsType[*LinkError](err)
				require.True(t, ok)

				require.Equal(t, tc.fieldsErrs, linkErr.Fields)
			}
		})
	}
}
