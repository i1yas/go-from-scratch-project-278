package httpapi

import (
	"fmt"

	"hexleturlshort/internal/links/application"
)

type contentRangeParams struct {
	itemName   string
	itemsRange application.Range
	itemsCount int
	totalCount int64
}

func formatContentRangeHeader(params contentRangeParams) string {
	r := params.itemsRange

	if params.itemsCount == 0 {
		return fmt.Sprintf("%s */%d",
			params.itemName,
			params.totalCount,
		)
	}

	return fmt.Sprintf("%s %d-%d/%d",
		params.itemName,
		r.From,
		max(0, r.From+int32(params.itemsCount)-1),
		params.totalCount,
	)
}
