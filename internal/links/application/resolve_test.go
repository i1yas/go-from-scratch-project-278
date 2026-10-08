package application_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
	"hexleturlshort/internal/links/application/testutils"
)

func TestResolveLink(t *testing.T) {
	t.Run("resolve existing link", func(t *testing.T) {
		linksStore := new(fakeLinksStore)
		visitsStore := new(fakeVisitsStore)

		params := commonResolveLinkParams("test")

		code, err := links.NewShortCode(params.Code)
		require.NoError(t, err)

		wantLink := testutils.Link(t, params.Code)

		linksStore.
			On("GetLinkByCode", mock.Anything, code).
			Return(wantLink, nil)

		visitsStore.
			On("CreateVisit", mock.Anything, mock.MatchedBy(func(v links.Visit) bool {
				return v.LinkID == wantLink.ID
			})).
			Return(nil)

		svc := application.NewService(
			linksStore,
			visitsStore,
			&fakeShortcodeGen{},
		)

		got, err := svc.ResolveLink(t.Context(), params)

		require.NoError(t, err)
		require.Equal(t, wantLink.OriginalURL, got)

		linksStore.AssertExpectations(t)
		visitsStore.AssertExpectations(t)
	})

	// TODO: check errors for visit recording: invalid record and store failure

	t.Run("resolve with invalid code", func(t *testing.T) {
		store := new(fakeLinksStore)

		params := commonResolveLinkParams("")

		svc := application.NewService(
			store,
			&fakeVisitsStore{},
			&fakeShortcodeGen{},
		)

		_, err := svc.ResolveLink(t.Context(), params)

		require.ErrorIs(t, err, links.ErrInvlalidShortCode)
		require.Equal(t, len(store.Calls), 0)
	})

	t.Run("link not found", func(t *testing.T) {
		linksStore := new(fakeLinksStore)

		params := commonResolveLinkParams("test")

		linksStore.
			On("GetLinkByCode", mock.Anything, mock.Anything).
			Return(links.Link{}, application.ErrLinkNotFound)

		svc := application.NewService(
			linksStore,
			&fakeVisitsStore{},
			&fakeShortcodeGen{},
		)

		_, err := svc.ResolveLink(t.Context(), params)

		require.ErrorIs(t, err, application.ErrLinkNotFound)
		require.Greater(t, len(linksStore.Calls), 0)
	})
}

func commonResolveLinkParams(codeRaw string) application.ResolveLinkParams {
	return application.ResolveLinkParams{
		Code:   codeRaw,
		IP:     "1.2.3.4",
		Status: 302,
	}
}
