package application

import "errors"

var (
	ErrLinkNotFound                    = errors.New("link not found")
	ErrShortCodeConflict               = errors.New("shortcode already exists")
	ErrFailedToGenerateUniqueShortCode = errors.New("failed to generate unique shortcode")
	ErrFailedToGenerateValidShortCode  = errors.New("failed to generate valid shortcode")
)
