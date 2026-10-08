package application_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links/application"
)

func TestDeleteLink(t *testing.T) {
	t.Run("delete existing link", func(t *testing.T) {
		store := new(fakeLinksStore)

		id := int64(101)

		store.
			On("DeleteLink", mock.Anything, id).
			Return(nil)

		svc := application.NewService(
			store,
			&fakeVisitsStore{},
			&fakeShortcodeGen{},
		)

		err := svc.DeleteLink(t.Context(), id)

		require.NoError(t, err)
		require.Equal(t, 1, len(store.Calls))
	})

	t.Run("error link not found", func(t *testing.T) {
		store := new(fakeLinksStore)

		store.
			On("DeleteLink", mock.Anything, mock.Anything).
			Return(application.ErrLinkNotFound)

		svc := application.NewService(
			store,
			&fakeVisitsStore{},
			&fakeShortcodeGen{},
		)

		err := svc.DeleteLink(t.Context(), 101)

		require.ErrorIs(t, err, application.ErrLinkNotFound)
		require.Equal(t, 1, len(store.Calls))
	})
}
