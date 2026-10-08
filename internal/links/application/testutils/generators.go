package testutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
)

// Link generates link based on shortcode
func Link(t *testing.T, codeRaw string) links.Link {
	t.Helper()

	link, err := links.NewLink(
		"http://test.com",
		codeRaw,
	)
	require.NoError(t, err)

	return link
}

// FullLink generates link based on all fields
func FullLink(
	t *testing.T,
	id int64,
	originalURL string,
	codeRaw string,
) links.Link {
	t.Helper()

	link, err := links.NewLinkWithID(
		id,
		originalURL,
		codeRaw,
	)
	require.NoError(t, err)

	return link
}

// Links generates links
func Links(t *testing.T, count int) []links.Link {
	result := make([]links.Link, count)

	for i := range count {
		link, err := links.NewLinkWithID(
			int64(i+1),
			fmt.Sprintf("http://test.com/page-%d", i),
			fmt.Sprintf("test-%d", i),
		)
		require.NoError(t, err)

		result[i] = link
	}

	return result
}

// Visits generates visits
func Visits(t *testing.T, count int) []links.Visit {
	result := make([]links.Visit, count)

	for i := range count {
		visit, err := links.NewVisit(
			1, "1.2.3.4", "", "", 302,
		)
		require.NoError(t, err)

		visit.ID = int64(i + 1)

		result[i] = visit
	}

	return result
}

// Range generates Range
func Range(t *testing.T, from, to int) application.Range {
	rang, err := application.NewRange(int32(from), int32(to))
	require.NoError(t, err)

	return rang
}

// SortOrder generates SortOrder
func SortOrder(t *testing.T, sortBy, order string) application.SortOrder {
	sort, err := application.NewSortOrder(sortBy, order)
	require.NoError(t, err)

	return sort
}
