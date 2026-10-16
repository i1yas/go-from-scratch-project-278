package application

import (
	"errors"
	"fmt"
	"math"
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
		return Range{}, fmt.Errorf("%w: negative lower bound", ErrInvalidRange)
	}

	if from > to {
		return Range{}, fmt.Errorf("%w: lower bound is greater that upper bound", ErrInvalidRange)
	}

	if to == math.MaxInt32 {
		return Range{}, fmt.Errorf("%w: upper bound is too big", ErrInvalidRange)
	}

	// NOTE: we can skip from == math.MaxInt32 check, because from > to already caught it

	return Range{From: from, To: to}, nil
}

// Size calculates amount of elements in range
func (r Range) Size() int32 {
	return r.To + 1 - r.From
}
