package httpapi

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"

	"hexleturlshort/internal/links/application"
)

func parseRangeFromReq(c *gin.Context, maxSize int32, defaultRange application.Range) (application.Range, error) {
	rawRange := c.Query("range")
	if rawRange == "" {
		rawRange = c.GetHeader("Range")
	}

	return parseRange(rawRange, maxSize, defaultRange)
}

func parseRange(raw string, maxSize int32, defaultRange application.Range) (application.Range, error) {
	if raw == "" {
		return defaultRange, nil
	}

	type rangeShape []int32

	var rangeItems rangeShape

	err := json.Unmarshal([]byte(raw), &rangeItems)
	if err != nil {
		return application.Range{}, fmt.Errorf("%w: invalid range format", ErrInvalidQuery)
	}

	if len(rangeItems) < 2 {
		return application.Range{}, fmt.Errorf("%w: invalid range format", ErrInvalidQuery)
	}

	rang, err := application.NewRange(rangeItems[0], rangeItems[1])
	if err != nil {
		return application.Range{}, err
	}

	if rang.Size() > maxSize {
		return application.Range{}, fmt.Errorf("%w: max size is %d", ErrRangeIsTooBig, maxSize)
	}

	return rang, nil
}

func parseSort(raw string, defaultSort application.SortOrder) (application.SortOrder, error) {
	if raw == "" {
		return defaultSort, nil
	}

	type sortShape []string

	var sortItems sortShape

	err := json.Unmarshal([]byte(raw), &sortItems)
	if err != nil {
		return application.SortOrder{}, fmt.Errorf("%w: invalid sort format", ErrInvalidQuery)
	}

	if len(sortItems) < 2 {
		return application.SortOrder{}, fmt.Errorf("%w: invalid sort format", ErrInvalidQuery)
	}

	return application.NewSortOrder(sortItems[0], sortItems[1])
}
