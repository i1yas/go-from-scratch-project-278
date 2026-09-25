package links

import "errors"

// ShortCode is object value that represents code for short url
type ShortCode string

var ErrInvlalidShortCode = errors.New("invalid shortcode")

// NewShortCode validates and creates valid shortcode
func NewShortCode(codeRaw string) (ShortCode, error) {
	l := len(codeRaw)
	if l < 3 || l > 32 {
		return ShortCode(""), ErrInvlalidShortCode
	}

	return ShortCode(codeRaw), nil
}
