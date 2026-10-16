package links

import (
	"errors"
	"fmt"
	"unicode"
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
	err := validateStringLength(codeRaw, shortcodeMinLength, shortcodeMaxLength)
	if err != nil {
		return ShortCode(""), fmt.Errorf("%w: %w",
			ErrInvlalidShortCode, err)
	}

	for _, r := range codeRaw {
		if unicode.IsLetter(r) || unicode.IsDigit(r) ||
			r == '-' || r == '_' || r == '.' || r == '@' {
			continue
		}

		return ShortCode(""), fmt.Errorf("%w: invalid character '%c' (%U)",
			ErrInvlalidShortCode, r, r,
		)
	}

	return ShortCode(codeRaw), nil
}
