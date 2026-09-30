package application

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// OrderASC is ASC order
	OrderASC = "ASC"
	// OrderDESC is DESC order
	OrderDESC = "DESC"
)

// SortOrder contains field and order for sorting list
type SortOrder struct {
	sortBy string
	order  string
}

var (
	ErrInvalidSortOrder     = errors.New("invalid sort order")
	ErrUnsupportedSortOrder = errors.New("unsupported sort order")
)

// NewSortOrder validates input and creates SortOrder
func NewSortOrder(sortBy, order string) (SortOrder, error) {
	if sortBy == "" {
		return SortOrder{}, fmt.Errorf("%w: empty sort field", ErrInvalidSortOrder)
	}

	switch order {
	case OrderASC, OrderDESC:
		break
	default:
		return SortOrder{}, fmt.Errorf("%w: invalid order", ErrInvalidSortOrder)
	}

	return SortOrder{
		sortBy: sortBy,
		order:  order,
	}, nil
}

func (so *SortOrder) String() string {
	str := fmt.Sprintf("%s %s", so.sortBy, so.order)
	return strings.TrimSpace(str)
}
