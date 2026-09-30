package application

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
)

func TestUpdateLink(t *testing.T) {
	t.Run("update link with valid data", func(t *testing.T) {
		store := new(fakeLinksStore)

		codeRaw := "test"
		want := createValidLink(t, codeRaw)
		want.ID = 101

		store.
			On("UpdateLink", mock.Anything, want).
			Return(want, nil)

		svc := NewService(store, &fakeShortcodeGen{})

		got, err := svc.UpdateLink(t.Context(), UpdateLinkParams{
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

		svc := NewService(store, &fakeShortcodeGen{})

		_, err := svc.UpdateLink(t.Context(), UpdateLinkParams{
			ID:          101,
			OriginalURL: "http://test.com",
			ShortCode:   "",
		})

		var linkErr *links.LinkError
		require.ErrorAs(t, err, &linkErr)
		require.Contains(t, linkErr.Fields, "shortcode")
		require.Equal(t, 0, len(store.Calls))
	})

	t.Run("update with invalid original url", func(t *testing.T) {
		store := new(fakeLinksStore)

		svc := NewService(store, &fakeShortcodeGen{})

		_, err := svc.UpdateLink(t.Context(), UpdateLinkParams{
			ID:          101,
			OriginalURL: "",
			ShortCode:   "code",
		})

		var linkErr *links.LinkError
		require.ErrorAs(t, err, &linkErr)
		require.Contains(t, linkErr.Fields, "original_url")
		require.Equal(t, 0, len(store.Calls))
	})

	t.Run("link to update not found", func(t *testing.T) {
		store := new(fakeLinksStore)

		validLink := createValidLink(t, "test")

		store.
			On("UpdateLink", mock.Anything, mock.Anything).
			Return(links.Link{}, ErrLinkNotFound)

		svc := NewService(store, &fakeShortcodeGen{})

		_, err := svc.UpdateLink(t.Context(), UpdateLinkParams{
			ID:          101,
			OriginalURL: string(validLink.OriginalURL),
			ShortCode:   string(validLink.ShortCode),
		})

		require.ErrorIs(t, err, ErrLinkNotFound)
		require.Equal(t, 1, len(store.Calls))
	})

	t.Run("failed to update with conflicting code", func(t *testing.T) {
		store := new(fakeLinksStore)

		want := createValidLink(t, "test")
		want.ID = 101

		store.
			On("UpdateLink", mock.Anything, want).
			Return(want, ErrShortCodeConflict)

		svc := NewService(store, &fakeShortcodeGen{})

		_, err := svc.UpdateLink(t.Context(), UpdateLinkParams{
			ID:          want.ID,
			OriginalURL: string(want.OriginalURL),
			ShortCode:   string(want.ShortCode),
		})

		require.ErrorIs(t, err, ErrShortCodeConflict)
		require.Equal(t, 1, len(store.Calls))
	})
}
