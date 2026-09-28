package application

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteLink(t *testing.T) {
	store := new(fakeLinksStore)

	store.
		On("DeleteLink", mock.Anything, mock.Anything).
		Return(nil)

	svc := NewService(store, &fakeShortcodeGen{})

	err := svc.DeleteLink(t.Context(), 101)

	require.NoError(t, err)
}

// TODO: test other cases
