package application_test

import (
	"context"
	"errors"

	"github.com/stretchr/testify/mock"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
)

type fakeLinksStore struct {
	mock.Mock
}

func (s *fakeLinksStore) GetLinks(ctx context.Context, params application.GetLinksParams) ([]links.Link, error) {
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

func (s *fakeVisitsStore) GetVisits(ctx context.Context, params application.GetVisitsParams) ([]links.Visit, error) {
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
