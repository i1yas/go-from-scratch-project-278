package application

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
)

func TestGetLinks(t *testing.T) {
	t.Run("one link", func(t *testing.T) {
		store := new(fakeLinksStore)

		link := createValidLink(t, "test")

		linkItems := []links.Link{link}

		store.
			On("GetLinks", mock.Anything).
			Return(linkItems, nil)

		generator := &fakeShortcodeGen{}
		svc := NewService(store, generator)

		got, err := svc.GetLinks(t.Context())

		require.NoError(t, err)
		require.Equal(t, len(linkItems), len(got))
		require.Equal(t, linkItems, got)
	})

	t.Run("no links", func(t *testing.T) {
		store := new(fakeLinksStore)

		want := []links.Link{}

		store.
			On("GetLinks", mock.Anything).
			Return([]links.Link{}, nil)

		generator := &fakeShortcodeGen{}
		svc := NewService(store, generator)

		got, err := svc.GetLinks(t.Context())

		require.NoError(t, err)
		require.Equal(t, len(want), len(got))
		require.Equal(t, want, got)
	})
}

func TestGetLinkByID(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		store := new(fakeLinksStore)

		want := createValidLink(t, "test")
		want.ID = 101

		store.
			On("GetLinkByID", mock.Anything, want.ID).
			Return(want, nil)

		generator := &fakeShortcodeGen{}
		svc := NewService(store, generator)

		got, err := svc.GetLinkByID(t.Context(), want.ID)

		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("not found", func(t *testing.T) {
		store := new(fakeLinksStore)

		store.
			On("GetLinkByID", mock.Anything, mock.Anything).
			Return(links.Link{}, ErrLinkNotFound)

		generator := &fakeShortcodeGen{}
		svc := NewService(store, generator)

		_, err := svc.GetLinkByID(t.Context(), 101)

		require.ErrorIs(t, err, ErrLinkNotFound)
	})
}
