package links

import (
	"errors"
	"fmt"
	"net/url"
)

// URL is object value for representing valid URL.
type URL string

const (
	urlMinLength = 8 // NOTE: min valid url is http://x
	urlMaxLength = 200
)

var ErrInvlalidURL = errors.New("invalid url")

// NewURL validates and creates valid URL, or returns error if URL is invalid.
func NewURL(urlRaw string) (URL, error) {
	err := validateStringLength(urlRaw, urlMinLength, urlMaxLength)
	if err != nil {
		return URL(""), fmt.Errorf("%w: %w", ErrInvlalidURL, err)
	}

	parsed, err := url.Parse(urlRaw)
	if err != nil {
		return URL(""), fmt.Errorf("%w: %w", ErrInvlalidURL, err)
	}

	switch parsed.Scheme {
	case "http", "https":
		break
	default:
		return URL(""), fmt.Errorf("%w: unsupported protocol '%s'", ErrInvlalidURL, parsed.Scheme)
	}

	if len(parsed.Host) == 0 {
		return URL(""), fmt.Errorf("%w: missing host", ErrInvlalidURL)
	}

	return URL(parsed.String()), nil
}
