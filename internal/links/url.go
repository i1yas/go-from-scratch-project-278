package links

import (
	"errors"
	"fmt"
	"net/url"
)

// URL is object value for representing valid URL.
type URL string

var ErrInvlalidURL = errors.New("invalid url")

// NewURL validates and creates valid URL, or returns error if URL is invalid.
func NewURL(urlRaw string) (URL, error) {
	if len(urlRaw) == 0 {
		return URL(""), ErrInvlalidURL
	}

	parsed, err := url.Parse(urlRaw)
	if err != nil {
		return URL(""), fmt.Errorf("%w: %w", ErrInvlalidURL, err)
	}

	switch parsed.Scheme {
	case "http", "https":
		break
	default:
		return URL(""), ErrInvlalidURL
	}

	if len(parsed.Host) == 0 {
		return URL(""), ErrInvlalidURL
	}

	return URL(parsed.String()), nil
}
