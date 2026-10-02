package links

import (
	"errors"
	"fmt"
	"net/netip"
	"time"
)

var (
	ErrInvalidVisit = errors.New("invalid visit")
	ErrInvalidIP    = errors.New("invalid IP address")
)

// Visit contains info about link visit
type Visit struct {
	ID        int64
	LinkID    int64
	CreatedAt time.Time
	IP        netip.Addr
	Referer   string
	UserAgent string
	Status    int
}

// NewVisit validates input and creates Visit
func NewVisit(
	linkID int64,
	ip string,
	referer string,
	userAgent string,
	status int,
) (Visit, error) {
	parsedIP, err := netip.ParseAddr(ip)
	if err != nil {
		return Visit{}, fmt.Errorf("%w: %w: %w", ErrInvalidVisit, ErrInvalidIP, err)
	}

	return Visit{
		LinkID:    linkID,
		IP:        parsedIP,
		Referer:   referer,
		UserAgent: userAgent,
		Status:    status,
	}, nil
}
