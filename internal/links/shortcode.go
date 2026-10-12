package links

import (
	"errors"
	"fmt"
)

// ShortCode is object value that represents code for short url
type ShortCode string

var ErrInvlalidShortCode = errors.New("invalid shortcode")

const (
	shortcodeMinLength = 3
	shortcodeMaxLength = 32
)

// NewShortCode validates and creates valid shortcode
func NewShortCode(codeRaw string) (ShortCode, error) {
	l := len(codeRaw)

	if l < shortcodeMinLength {
		return ShortCode(""), fmt.Errorf("%w: too short, min length is %d",
			ErrInvlalidShortCode, shortcodeMinLength)
	}

	if l > shortcodeMaxLength {
		return ShortCode(""), fmt.Errorf("%w: too long, max length is %d",
			ErrInvlalidShortCode, shortcodeMaxLength)
	}

	return ShortCode(codeRaw), nil
}
