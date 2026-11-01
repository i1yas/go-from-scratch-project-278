package links

// Link represents information about link entry
type Link struct {
	ID          int64
	OriginalURL URL
	ShortCode   ShortCode
}

// NewLink takes raw input, validates and creates Link
func NewLink(originalURLRaw, shortcodeRaw string) (Link, error) {
	originalURL, err := NewURL(originalURLRaw)
	if err != nil {
		return Link{}, err
	}

	shortcode, err := NewShortCode(shortcodeRaw)
	if err != nil {
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
