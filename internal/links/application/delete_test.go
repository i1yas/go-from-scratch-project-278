package application

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteLink(t *testing.T) {
	t.Run("delete existing link", func(t *testing.T) {
		store := new(fakeLinksStore)

		id := int64(101)

		store.
			On("DeleteLink", mock.Anything, id).
			Return(nil)

		svc := NewService(store, &fakeShortcodeGen{})

		err := svc.DeleteLink(t.Context(), id)

		require.NoError(t, err)
		require.Equal(t, 1, len(store.Calls))
	})

	t.Run("error link not found", func(t *testing.T) {
		store := new(fakeLinksStore)

		store.
			On("DeleteLink", mock.Anything, mock.Anything).
			Return(ErrLinkNotFound)

		svc := NewService(store, &fakeShortcodeGen{})

		err := svc.DeleteLink(t.Context(), 101)

		require.ErrorIs(t, err, ErrLinkNotFound)
		require.Equal(t, 1, len(store.Calls))
	})
}
