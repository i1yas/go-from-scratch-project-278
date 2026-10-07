package application

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
)

type fakeLinksStore struct {
	mock.Mock
}

func (s *fakeLinksStore) GetLinks(ctx context.Context, params GetLinksParams) ([]links.Link, error) {
	args := s.Called(ctx, params)
	return args.Get(0).([]links.Link), args.Error(1)
}

func (s *fakeLinksStore) GetLinksTotalCount(ctx context.Context) (int64, error) {
	args := s.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (s *fakeLinksStore) GetLinkByID(ctx context.Context, id int64) (links.Link, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(links.Link), args.Error(1)
}

func (s *fakeLinksStore) GetLinkByCode(ctx context.Context, code links.ShortCode) (links.Link, error) {
	args := s.Called(ctx, code)
	return args.Get(0).(links.Link), args.Error(1)
}

func (s *fakeLinksStore) CreateLink(ctx context.Context, link links.Link) (links.Link, error) {
	args := s.Called(ctx, link)
	return args.Get(0).(links.Link), args.Error(1)
}

func (s *fakeLinksStore) UpdateLink(ctx context.Context, link links.Link) (links.Link, error) {
	args := s.Called(ctx, link)
	return args.Get(0).(links.Link), args.Error(1)
}

func (s *fakeLinksStore) DeleteLink(ctx context.Context, id int64) error {
	args := s.Called(ctx, id)
	return args.Error(0)
}

type fakeVisitsStore struct {
	mock.Mock
}

func (s *fakeVisitsStore) GetVisits(ctx context.Context, params GetVisitsParams) ([]links.Visit, error) {
	args := s.Called(ctx, params)
	return args.Get(0).([]links.Visit), args.Error(1)
}

func (s *fakeVisitsStore) GetVisitsTotalCount(ctx context.Context) (int64, error) {
	args := s.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (s *fakeVisitsStore) CreateVisit(ctx context.Context, link links.Visit) error {
	args := s.Called(ctx, link)
	return args.Error(0)
}

type fakeShortcodeGen struct {
	codes []string
	calls int
}

func (g *fakeShortcodeGen) Generate() (string, error) {
	ind := g.calls
	g.calls++

	if len(g.codes) == 0 {
		return "", errors.New("failed to generate")
	}

	code := g.codes[ind%len(g.codes)]

	return code, nil
}

func createValidLink(t *testing.T, codeRaw string) links.Link {
	t.Helper()

	link, err := links.NewLink(
		"http://test.com",
		codeRaw,
	)
	require.NoError(t, err)

	return link
}

func createVisits(t *testing.T, count int) []links.Visit {
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

func commonResolveLinkParams(codeRaw string) ResolveLinkParams {
	return ResolveLinkParams{
		Code:   codeRaw,
		IP:     "1.2.3.4",
		Status: 302,
	}
}

func createRange(t *testing.T, from, to int) Range {
	rang, err := NewRange(int32(from), int32(to))
	require.NoError(t, err)

	return rang
}

func createSort(t *testing.T, sortBy, order string) SortOrder {
	sort, err := NewSortOrder(sortBy, order)
	require.NoError(t, err)

	return sort
}
