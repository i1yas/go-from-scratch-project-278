package application

import "errors"

var (
	ErrLinkNotFound                    = errors.New("link not found")
	ErrShortCodeConflict               = errors.New("shortcode already exists")
	ErrFailedToGenerateUniqueShortCode = errors.New("failed to generate unique shortcode")
	ErrShortCodeGeneratorInternal      = errors.New("shortcode generator internal error")
	ErrStoreInternal                   = errors.New("store internal error")
	ErrInvalidStoreValue               = errors.New("invalid store value")
)
