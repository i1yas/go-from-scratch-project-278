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
		errs        []FieldError
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
			errs: []FieldError{
				{"original_url", ErrInvlalidURL},
			},
		},
		{
			name:        "error, empty shortcode",
			originalURL: "http://test.com",
			shortcode:   "",
			errs: []FieldError{
				{"shortcode", ErrInvlalidShortCode},
			},
		},
		{
			name:        "error, both original url and shortcode are empty",
			originalURL: "",
			shortcode:   "",
			errs: []FieldError{
				{"original_url", ErrInvlalidURL},
				{"shortcode", ErrInvlalidShortCode},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewLink(tc.originalURL, tc.shortcode)

			if tc.errs == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				var verr *ValidationError
				require.ErrorAs(t, err, &verr)
				require.Equal(t, len(tc.errs), len(verr.Fields))

				for _, ferr := range tc.errs {
					require.ErrorIs(t, err, ferr.Err)
					require.ErrorContains(t, err, ferr.Field)
				}
			}
		})
	}
}
