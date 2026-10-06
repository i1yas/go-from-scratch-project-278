package postgres

import (
	"database/sql"
	"net"
	"net/netip"

	"github.com/sqlc-dev/pqtype"

	"hexleturlshort/internal/links/application"
)

func convertRangeToLimitOffset(rang application.Range) (int32, int32) {
	// NOTE: upper bound is inclusive
	limit := rang.To + 1 - rang.From

	offset := rang.From
	if rang.From == rang.To {
		limit = 0
	}

	return limit, offset
}

func nullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func addrToPgInet(addr netip.Addr) pqtype.Inet {
	return pqtype.Inet{
		IPNet: net.IPNet{
			IP:   net.IP(addr.AsSlice()),
			Mask: net.IPMask{0, 0, 0, 0},
		},
		Valid: true,
	}
}
