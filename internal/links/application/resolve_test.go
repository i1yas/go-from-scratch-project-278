package application

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
)

func TestResolveLink(t *testing.T) {
	t.Run("resolve existing link", func(t *testing.T) {
		store := new(fakeLinksStore)

		codeRaw := "test"

		code, err := links.NewShortCode(codeRaw)
		require.NoError(t, err)

		want := createValidLink(t, codeRaw)

		store.
			On("GetLinkByCode", mock.Anything, code).
			Return(want, nil)

		svc := NewService(store, &fakeShortcodeGen{})

		got, err := svc.ResolveLink(t.Context(), codeRaw)

		require.NoError(t, err)
		require.Equal(t, want.OriginalURL, got)
		require.Greater(t, len(store.Calls), 0)
	})

	t.Run("resolve with invalid code", func(t *testing.T) {
		store := new(fakeLinksStore)

		emptyCode := ""

		svc := NewService(store, &fakeShortcodeGen{})

		_, err := svc.ResolveLink(t.Context(), emptyCode)

		require.ErrorIs(t, err, links.ErrInvlalidShortCode)
		require.Equal(t, len(store.Calls), 0)
	})

	t.Run("link not found", func(t *testing.T) {
		store := new(fakeLinksStore)

		codeRaw := "test"

		store.
			On("GetLinkByCode", mock.Anything, mock.Anything).
			Return(links.Link{}, ErrLinkNotFound)

		svc := NewService(store, &fakeShortcodeGen{})

		_, err := svc.ResolveLink(t.Context(), codeRaw)

		require.ErrorIs(t, err, ErrLinkNotFound)
		require.Greater(t, len(store.Calls), 0)
	})
}
