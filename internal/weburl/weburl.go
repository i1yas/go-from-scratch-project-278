package weburl

import (
	"errors"
	"fmt"
	"net/netip"
	"net/url"

	"golang.org/x/net/idna"
)

var (
	ErrFailedToParse        = errors.New("failed to parse")
	ErrMissingHost          = errors.New("missing host")
	ErrUnsupportedProtocol  = errors.New("unsupported protocol")
	ErrIDNAConversionFailed = errors.New("IDNA conversion failed")
)

var idnaProfile = idna.New(
	// NOTE: convert alternative unicode characters, reject unsupported characters (UTS #46, IETF/STD 3)
	idna.MapForLookup(),
	idna.Transitional(false),
	// NOTE: validate bidirectional scripts, applied to text that contains characters from right-to-left scripts
	idna.BidiRule(),
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

	if _, err := netip.ParseAddr(baseURL.Hostname()); err == nil {
		return nil, fmt.Errorf("%w: ip address is not allowed", ErrFailedToParse)
	}

	// NOTE: validate, normalize and convert to punycode
	hostname, err := idnaProfile.ToASCII(baseURL.Hostname())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrIDNAConversionFailed, err)
	}

	host := hostname
	if baseURL.Port() != "" {
		host += ":" + baseURL.Port()
	}

	baseURL.Host = host

	return baseURL, nil
}
