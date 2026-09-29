package links

import "net/url"

// Link represents information about link entry
type Link struct {
	ID          int64
	OriginalURL URL
	ShortCode   ShortCode
}

// NewLink takes raw input, validates and creates Link
func NewLink(originalURLRaw, shortcodeRaw string) (Link, error) {
	fieldsErrs := make(map[string]string)

	originalURL, err := NewURL(originalURLRaw)
	if err != nil {
		fieldsErrs["original_url"] = err.Error()
	}

	shortcode, err := NewShortCode(shortcodeRaw)
	if err != nil {
		fieldsErrs["shortcode"] = err.Error()
	}

	if len(fieldsErrs) > 0 {
		return Link{}, &LinkError{
			Fields: fieldsErrs,
		}
	}

	link := Link{
		OriginalURL: originalURL,
		ShortCode:   shortcode,
	}

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

// LinkError contains field Link field errors
type LinkError struct {
	Fields map[string]string
}

// Error returns LinkError message
func (e *LinkError) Error() string {
	return "invalid link"
}
