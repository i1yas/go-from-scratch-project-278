package application

import (
	"errors"
	"fmt"
)

var ErrInvalidRange = errors.New("invalid range")

// Range represents range for lists
type Range struct {
	From int32
	To   int32
}

// NewRange validates bounds and creates Range
func NewRange(from, to int32) (Range, error) {
	if from < 0 {
		return Range{}, fmt.Errorf("%w: negative 'from'", ErrInvalidRange)
	}

	if to < 0 {
		return Range{}, fmt.Errorf("%w: negative 'to'", ErrInvalidRange)
	}

	if from > to {
		return Range{}, fmt.Errorf("%w: 'from' is bigger than 'to'", ErrInvalidRange)
	}

	return Range{From: from, To: to}, nil
}
