package links

import "net/url"

// Link represents information about link entry
type Link struct {
	ID          int64
	OriginalURL URL
	ShortCode   ShortCode
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
