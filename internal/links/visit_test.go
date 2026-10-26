package links

import (
	"net/netip"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewVisit(t *testing.T) {
	longURL := "http://" + strings.Repeat("x", 1000)
	longURLWithCyrillicHost := "http://пример.рф" + strings.Repeat("x", 1000)

	cases := []struct {
		name      string
		ip        string
		referer   string
		userAgent string
		status    int
		want      Visit
		errs      []error
	}{
		{
			name:    "valid visit with IP v4",
			ip:      "10.11.12.13",
			referer: "http://mysite",
			status:  302,
			want: Visit{
				IP:      netip.AddrFrom4([4]byte{10, 11, 12, 13}),
				Referer: "http://mysite",
				Status:  302,
			},
		},
		{
			name:    "valid visit with IP v6 full-form",
			ip:      "2001:db8:3333:4444:5555:6666:7777:8888",
			referer: "http://mysite",
			status:  302,
			want: Visit{
				IP: netip.AddrFrom16([16]byte{
					32, 1, 13, 184,
					51, 51, 68, 68,
					85, 85, 102, 102,
					119, 119, 136, 136,
				}),
				Referer: "http://mysite",
				Status:  302,
			},
		},
		{
			name:    "valid visit with IP v6 skip segements",
			ip:      "2001:db8::",
			referer: "http://mysite",
			status:  302,
			want: Visit{
				IP: netip.AddrFrom16([16]byte{
					32, 1, 13, 184,
				}),
				Referer: "http://mysite",
				Status:  302,
			},
		},
		{
			name:    "valid visit with max len referer",
			ip:      "10.11.12.13",
			referer: longURL[0:refererMaxLength],
			status:  302,
			want: Visit{
				IP:      netip.AddrFrom4([4]byte{10, 11, 12, 13}),
				Referer: longURL[0:refererMaxLength],
				Status:  302,
			},
		},
		{
			name:    "valid visit with max len unicode referer",
			ip:      "10.11.12.13",
			referer: longURLWithCyrillicHost[0:refererMaxLength],
			status:  302,
			want: Visit{
				IP:      netip.AddrFrom4([4]byte{10, 11, 12, 13}),
				Referer: longURLWithCyrillicHost[0:refererMaxLength],
				Status:  302,
			},
		},
		{
			name:      "valid visit with max len user-agent",
			ip:        "10.11.12.13",
			referer:   "http://mysite",
			userAgent: strings.Repeat("x", userAgentMaxLength),
			status:    302,
			want: Visit{
				IP:        netip.AddrFrom4([4]byte{10, 11, 12, 13}),
				Referer:   "http://mysite",
				UserAgent: strings.Repeat("x", userAgentMaxLength),
				Status:    302,
			},
		},
		{
			name:      "valid visit with max len unicode user-agent",
			ip:        "10.11.12.13",
			referer:   "http://mysite",
			userAgent: strings.Repeat("щ", userAgentMaxLength),
			status:    302,
			want: Visit{
				IP:        netip.AddrFrom4([4]byte{10, 11, 12, 13}),
				Referer:   "http://mysite",
				UserAgent: strings.Repeat("щ", userAgentMaxLength),
				Status:    302,
			},
		},
		{
			name:    "invalid ip",
			ip:      "1.1",
			referer: "http://mysite",
			status:  302,
			errs:    []error{ErrInvalidVisit, ErrInvalidIP},
		},
		{
			name:    "invalid referer",
			ip:      "10.11.12.13",
			referer: longURL[0 : refererMaxLength+1],
			status:  302,
			errs:    []error{ErrInvalidVisit, ErrInvalidReferer},
		},
		{
			name:      "invalid user-agent",
			ip:        "10.11.12.13",
			referer:   "http://mysite",
			userAgent: strings.Repeat("x", userAgentMaxLength+1),
			status:    302,
			errs:      []error{ErrInvalidVisit, ErrInvalidUserAgent},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewVisit(
				0,
				tc.ip,
				tc.referer,
				tc.userAgent,
				tc.status,
			)

			for _, wantErr := range tc.errs {
				require.ErrorIs(t, err, wantErr)
			}

			require.Equal(t, tc.want, got)
		})
	}
}
