package links

import (
	"errors"
	"fmt"
	"net/netip"
	"strings"
	"time"
)

var (
	ErrInvalidVisit     = errors.New("invalid visit")
	ErrInvalidIP        = errors.New("invalid IP address")
	ErrInvalidReferer   = errors.New("invalid referer")
	ErrInvalidUserAgent = errors.New("invalid user-agent")
)

const (
	refererMaxLength   = 255
	userAgentMaxLength = 255
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

	userAgent = truncate(userAgent, userAgentMaxLength)
	userAgent = strings.ToValidUTF8(userAgent, "\uFFFD")

	referer = truncate(referer, refererMaxLength)
	referer = strings.ToValidUTF8(referer, "\uFFFD")

	return Visit{
		LinkID:    linkID,
		IP:        parsedIP,
		Referer:   truncate(referer, refererMaxLength),
		UserAgent: userAgent,
		Status:    status,
	}, nil
}

func truncate(s string, maxLen int) string {
	// NOTE: count != i, count - character count, i - byte index of character
	count := 0
	for i := range s {
		if count == maxLen {
			return s[:i]
		}

		count++
	}

	return s
}
