package weburl

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var (
	ErrFailedToParse       = errors.New("failed to parse")
	ErrMissingHost         = errors.New("missing host")
	ErrUnsupportedProtocol = errors.New("unsupported protocol")
)

// Parse parses raw url, validates and normalizes it
func Parse(rawURL string) (*url.URL, error) {
	baseURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToParse, err)
	}

	switch baseURL.Scheme {
	case "http", "https":
		break
	default:
		return nil, fmt.Errorf("%w '%s'", ErrUnsupportedProtocol, baseURL.Scheme)
	}

	if baseURL.Hostname() == "" {
		return nil, ErrMissingHost
	}

	baseURL.Host = strings.ToLower(baseURL.Host)

	return baseURL, nil
}
