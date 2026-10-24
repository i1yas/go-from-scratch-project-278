package httpapi

import (
	"errors"
	"fmt"
	"net/url"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/weburl"
)

var (
	ErrInvalidBaseURL           = errors.New("invalid base url")
	ErrFailedToGenerateShortURL = errors.New("failed to generate short url")
)

// ShortURLGenerator contains validated base url and generates short urls based on it
type ShortURLGenerator struct {
	baseURL *url.URL
}

// NewShortURLGenerator validates base url and creates short url generator
func NewShortURLGenerator(baseURLRaw string) (*ShortURLGenerator, error) {
	baseURL, err := weburl.Parse(baseURLRaw)
	if err != nil {
		return nil, fmt.Errorf("%w: %w (baseURL=%s)", ErrInvalidBaseURL, err, baseURLRaw)
	}

	return &ShortURLGenerator{
		baseURL: baseURL,
	}, nil
}

// Generate generates short urls
func (g *ShortURLGenerator) Generate(code links.ShortCode) (string, error) {
	shortURL := *g.baseURL

	path, err := url.JoinPath("r", string(code))
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToGenerateShortURL, err)
	}

	shortURL.Path = path

	return shortURL.String(), nil
}
