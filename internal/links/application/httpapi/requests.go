package httpapi

import (
	"encoding/json"
	"hexleturlshort/internal/links/application"
)

func parseRange(raw string, defaultRange application.Range) (application.Range, error) {
	if raw == "" {
		return defaultRange, nil
	}

	type rangeShape []int32

	var rangeItems rangeShape

	err := json.Unmarshal([]byte(raw), &rangeItems)
	if err != nil {
		return application.Range{}, ErrInvalidJSON
	}

	if len(rangeItems) < 2 {
		return defaultRange, nil
	}

	return application.NewRange(rangeItems[0], rangeItems[1])
}

func parseSort(raw string, defaultSort application.SortOrder) (application.SortOrder, error) {
	if raw == "" {
		return defaultSort, nil
	}

	type sortShape []string

	var sortItems sortShape

	err := json.Unmarshal([]byte(raw), &sortItems)
	if err != nil {
		return application.SortOrder{}, ErrInvalidJSON
	}

	if len(sortItems) < 2 {
		return defaultSort, nil
	}

	return application.NewSortOrder(sortItems[0], sortItems[1])
}
