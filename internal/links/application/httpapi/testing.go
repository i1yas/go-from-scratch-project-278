package httpapi

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
)

type fakeLinkService struct {
	mock.Mock
}

func (s *fakeLinkService) GetLinks(
	ctx context.Context,
	params application.GetLinksParams,
) (application.LinksResult, error) {
	args := s.Called(ctx, params)
	return args.Get(0).(application.LinksResult), args.Error(1)
}

func (s *fakeLinkService) GetLinkByID(ctx context.Context, id int64) (links.Link, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(links.Link), args.Error(1)
}

func (s *fakeLinkService) CreateLink(ctx context.Context, params application.CreateLinkParams) (links.Link, error) {
	args := s.Called(ctx, params)
	return args.Get(0).(links.Link), args.Error(1)
}

func (s *fakeLinkService) UpdateLink(ctx context.Context, params application.UpdateLinkParams) (links.Link, error) {
	args := s.Called(ctx, params)
	return args.Get(0).(links.Link), args.Error(1)
}

func (s *fakeLinkService) DeleteLink(ctx context.Context, id int64) error {
	args := s.Called(ctx, id)
	return args.Error(0)
}

func (s *fakeLinkService) ResolveLink(ctx context.Context, params application.ResolveLinkParams) (links.URL, error) {
	args := s.Called(ctx, params.Code)
	return args.Get(0).(links.URL), args.Error(1)
}

type fakeVisitsService struct {
	mock.Mock
}

func (s *fakeVisitsService) GetVisits(
	ctx context.Context, params application.GetVisitsParams,
) (application.VisitsResult, error) {
	args := s.Called(ctx, params)
	return args.Get(0).(application.VisitsResult), args.Error(1)
}

func testingBaseURL(t *testing.T) links.URL {
	baseURL, err := links.NewURL("http://short")
	require.NoError(t, err)

	return baseURL
}

func setupTestRouter(
	t *testing.T,
	linksSvc *fakeLinkService,
	visitsSvc *fakeVisitsService,
) *gin.Engine {
	baseURL := testingBaseURL(t)
	linksHandler := NewLinksHandler(linksSvc, baseURL)
	visitsHandler := NewVisitsHandler(visitsSvc, baseURL)
	router := gin.New()

	RegisterRoutes(router, linksHandler, visitsHandler)

	return router
}

func createValidLink(t *testing.T, id int64, originalURL, code string) links.Link {
	link, err := links.NewLinkWithID(id, originalURL, code)
	require.NoError(t, err)

	return link
}

func createLinks(t *testing.T, ids []int64) []links.Link {
	result := make([]links.Link, len(ids))

	for i, id := range ids {
		link, err := links.NewLinkWithID(
			id,
			fmt.Sprintf("http://domain.com/%d", id),
			fmt.Sprintf("code-%d", 100+id),
		)

		require.NoError(t, err)

		result[i] = link
	}

	return result
}

func createVisits(t *testing.T, ids []int64) []links.Visit {
	result := make([]links.Visit, len(ids))

	timeOffset, err := time.Parse(time.RFC3339, "2026-06-10T12:00:00Z")
	require.NoError(t, err)

	for i, id := range ids {
		ip := fmt.Sprintf("192.168.0.%d", id)
		if id == 10 {
			ip = "255.0.0.1"
		}

		visit, err := links.NewVisit(
			1,
			ip,
			fmt.Sprintf("http://website.com/page-%d", id),
			fmt.Sprintf("user-agent-%d", id),
			302,
		)
		require.NoError(t, err)

		visit.ID = id
		visit.CreatedAt = timeOffset.Add(
			time.Minute * time.Duration(id),
		)

		result[i] = visit
	}

	return result
}

func createRange(t *testing.T, from, to int) application.Range {
	t.Helper()

	r, err := application.NewRange(int32(from), int32(to))
	require.NoError(t, err)

	return r
}

func createSort(t *testing.T, sortBy, order string) application.SortOrder {
	t.Helper()

	sort, err := application.NewSortOrder(sortBy, order)
	require.NoError(t, err)

	return sort
}

func loadFixture(t *testing.T, path string) string {
	t.Helper()

	fixturePath, err := filepath.Abs(filepath.Join("testdata", "fixture", path))
	require.NoError(t, err)

	data, err := os.ReadFile(fixturePath)
	require.NoError(t, err)

	return string(data)
}

func loadResponseFixture(t *testing.T, name string) string {
	path := filepath.Join("responses", name+".json")

	return loadFixture(t, path)
}
