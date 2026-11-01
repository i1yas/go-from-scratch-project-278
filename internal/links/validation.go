package links

import (
	"fmt"
	"unicode/utf8"
)

func validateCharactersCount(s string, minLength, maxLength int) error {
	l := utf8.RuneCountInString(s)

	if l < minLength {
		return fmt.Errorf("too short, min length is %d", minLength)
	}

	if l > maxLength {
		return fmt.Errorf("too long, max length is %d", maxLength)
	}

	return nil
}
