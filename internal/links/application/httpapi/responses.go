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

	return fmt.Sprintf("%s %d-%d/%d",
		params.itemName,
		r.From,
		r.From+int32(params.itemsCount)-1,
		params.totalCount,
	)
}
