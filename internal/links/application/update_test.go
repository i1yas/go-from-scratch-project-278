package application

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateLink(t *testing.T) {
	store := new(fakeLinksStore)

	codeRaw := "test"
	want := createValidLink(t, codeRaw)
	want.ID = 101

	store.
		On("UpdateLink", mock.Anything, mock.Anything).
		Return(want, nil)

	svc := NewService(store, &fakeShortcodeGen{})

	got, err := svc.UpdateLink(t.Context(), UpdateLinkParams{
		ID:          want.ID,
		OriginalURL: string(want.OriginalURL),
		ShortCode:   codeRaw,
	})

	require.NoError(t, err)
	require.Equal(t, want, got)
}

// TODO: test other cases
