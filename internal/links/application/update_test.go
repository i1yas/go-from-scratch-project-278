package application_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
	"hexleturlshort/internal/links/application/testutils"
)

func TestUpdateLink(t *testing.T) {
	t.Run("update link with valid data", func(t *testing.T) {
		store := new(fakeLinksStore)

		codeRaw := "test"
		want := testutils.Link(t, codeRaw)
		want.ID = 101

		store.
			On("UpdateLink", mock.Anything, want).
			Return(want, nil)

		svc := application.NewService(
			store,
			&fakeVisitsStore{},
			&fakeShortcodeGen{},
		)

		got, err := svc.UpdateLink(t.Context(), application.UpdateLinkParams{
			ID:          want.ID,
			OriginalURL: string(want.OriginalURL),
			ShortCode:   codeRaw,
		})

		require.NoError(t, err)
		require.Equal(t, want, got)
		require.Equal(t, 1, len(store.Calls))
	})

	t.Run("update with invalid shortcode", func(t *testing.T) {
		store := new(fakeLinksStore)

		svc := application.NewService(
			store,
			&fakeVisitsStore{},
			&fakeShortcodeGen{},
		)

		_, err := svc.UpdateLink(t.Context(), application.UpdateLinkParams{
			ID:          101,
			OriginalURL: "http://test.com",
			ShortCode:   "",
		})

		require.ErrorContains(t, err, "shortcode:")
		require.ErrorIs(t, err, links.ErrInvlalidShortCode)
		require.Equal(t, 0, len(store.Calls))
	})

	t.Run("update with invalid original url", func(t *testing.T) {
		store := new(fakeLinksStore)

		svc := application.NewService(
			store,
			&fakeVisitsStore{},
			&fakeShortcodeGen{},
		)

		_, err := svc.UpdateLink(t.Context(), application.UpdateLinkParams{
			ID:          101,
			OriginalURL: "",
			ShortCode:   "code",
		})

		require.ErrorContains(t, err, "original_url:")
		require.ErrorIs(t, err, links.ErrInvlalidURL)
		require.Equal(t, 0, len(store.Calls))
	})

	t.Run("link to update not found", func(t *testing.T) {
		store := new(fakeLinksStore)

		validLink := testutils.Link(t, "test")

		store.
			On("UpdateLink", mock.Anything, mock.Anything).
			Return(links.Link{}, application.ErrLinkNotFound)

		svc := application.NewService(
			store,
			&fakeVisitsStore{},
			&fakeShortcodeGen{},
		)

		_, err := svc.UpdateLink(t.Context(), application.UpdateLinkParams{
			ID:          101,
			OriginalURL: string(validLink.OriginalURL),
			ShortCode:   string(validLink.ShortCode),
		})

		require.ErrorIs(t, err, application.ErrLinkNotFound)
		require.Equal(t, 1, len(store.Calls))
	})

	t.Run("failed to update with conflicting code", func(t *testing.T) {
		store := new(fakeLinksStore)

		want := testutils.Link(t, "test")
		want.ID = 101

		store.
			On("UpdateLink", mock.Anything, want).
			Return(want, application.ErrShortCodeConflict)

		svc := application.NewService(
			store,
			&fakeVisitsStore{},
			&fakeShortcodeGen{},
		)

		_, err := svc.UpdateLink(t.Context(), application.UpdateLinkParams{
			ID:          want.ID,
			OriginalURL: string(want.OriginalURL),
			ShortCode:   string(want.ShortCode),
		})

		require.ErrorIs(t, err, application.ErrShortCodeConflict)
		require.Equal(t, 1, len(store.Calls))
	})
}
