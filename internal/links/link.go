package links

import (
	"net/url"
)

// Link represents information about link entry
type Link struct {
	ID          int64
	OriginalURL URL
	ShortCode   ShortCode
}

// NewLink takes raw input, validates and creates Link
func NewLink(originalURLRaw, shortcodeRaw string) (Link, error) {
	var verr ValidationError

	originalURL, err := NewURL(originalURLRaw)
	verr.add("original_url", err)

	shortcode, err := NewShortCode(shortcodeRaw)
	verr.add("shortcode", err)

	if err := verr.err(); err != nil {
		return Link{}, err
	}

	link := Link{
		OriginalURL: originalURL,
		ShortCode:   shortcode,
	}

	return link, nil
}

// NewLinkWithID same as NewLink but with id
func NewLinkWithID(id int64, originalURLRaw, shortcodeRaw string) (Link, error) {
	link, err := NewLink(originalURLRaw, shortcodeRaw)
	if err != nil {
		return Link{}, err
	}

	link.ID = id

	return link, nil
}

// ShortURL generates url based on shortcode and base url
func (l Link) ShortURL(baseURL URL) (URL, error) {
	parsed, err := url.Parse(string(baseURL))
	if err != nil {
		return URL(""), err
	}

	path, err := url.JoinPath("r", string(l.ShortCode))
	if err != nil {
		return URL(""), err
	}

	parsed.Path = path

	return NewURL(parsed.String())
}
