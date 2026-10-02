package links

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewVisit(t *testing.T) {
	cases := []struct {
		name      string
		ip        string
		referer   string
		userAgent string
		status    int
		want      Visit
		err       error
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

			require.ErrorIs(t, err, tc.err)
			require.Equal(t, tc.want, got)
		})
	}
}
