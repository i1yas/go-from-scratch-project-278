package application

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResolveLink(t *testing.T) {
	store := new(fakeLinksStore)

	codeRaw := "test"
	want := createValidLink(t, codeRaw)

	store.
		On("GetLinkByCode", mock.Anything, mock.Anything).
		Return(want, nil)

	svc := NewService(store, &fakeShortcodeGen{})

	got, err := svc.ResolveLink(t.Context(), "codeRaw")

	require.NoError(t, err)
	require.Equal(t, want.OriginalURL, got)
}

// TODO: test other cases
