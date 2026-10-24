package links

import (
	"errors"
	"fmt"

	"hexleturlshort/internal/weburl"
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
	parsed, err := weburl.Parse(urlRaw)
	if err != nil {
		return URL(""), fmt.Errorf("%w: %w", ErrInvlalidURL, err)
	}

	normalizedURL := parsed.String()

	err = validateStringLength(normalizedURL, urlMinLength, urlMaxLength)
	if err != nil {
		return URL(""), fmt.Errorf("%w: %w", ErrInvlalidURL, err)
	}

	return URL(normalizedURL), nil
}
