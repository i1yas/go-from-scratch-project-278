package application

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
)

func TestCreateLink(t *testing.T) {
	t.Run("create with valid shortcode", func(t *testing.T) {
		store := new(fakeLinksStore)

		codeRaw := "test"
		params := CreateLinkParams{
			OriginalURL: "https://test.com/test",
			ShortCode:   &codeRaw,
		}

		want := links.Link{
			ID:          1,
			OriginalURL: links.URL(params.OriginalURL),
			ShortCode:   links.ShortCode(*params.ShortCode),
		}

		store.
			On("CreateLink", mock.Anything, mock.Anything).
			Return(want, nil)

		generator := &fakeShortcodeGen{}
		svc := NewService(store, generator)

		got, err := svc.CreateLink(t.Context(), params)

		require.NoError(t, err)
		require.Equal(t, want, got)
		require.Equal(t, generator.calls, 0)
	})

	t.Run("create without shortcode", func(t *testing.T) {
		store := new(fakeLinksStore)

		generatedCode := "generated"

		params := CreateLinkParams{
			OriginalURL: "https://test.com/test",
		}

		want := links.Link{
			ID:          1,
			OriginalURL: links.URL(params.OriginalURL),
			ShortCode:   links.ShortCode(generatedCode),
		}

		store.
			On("CreateLink", mock.Anything, mock.Anything).
			Return(want, nil)

		generator := &fakeShortcodeGen{codes: []string{generatedCode}}
		svc := NewService(store, generator)

		got, err := svc.CreateLink(t.Context(), params)

		require.NoError(t, err)
		require.Equal(t, want, got)
		require.Equal(t, generator.calls, 1)
	})

	t.Run("create without shortcode with retry", func(t *testing.T) {
		store := new(fakeLinksStore)

		params := CreateLinkParams{
			OriginalURL: "https://test.com/test",
		}

		want := links.Link{
			ID:          1,
			OriginalURL: links.URL(params.OriginalURL),
			ShortCode:   links.ShortCode("code-two"),
		}

		store.
			On("CreateLink", mock.Anything, mock.Anything).
			Return(links.Link{}, ErrShortCodeConflict).
			Once().
			On("CreateLink", mock.Anything, mock.Anything).
			Return(want, nil)

		generator := &fakeShortcodeGen{codes: []string{"code-one", "code-two"}}
		svc := NewService(store, generator)

		got, err := svc.CreateLink(t.Context(), params)

		require.NoError(t, err)
		require.Equal(t, want, got)
		require.Equal(t, generator.calls, 2)
	})

	t.Run("create with invalid shortcode", func(t *testing.T) {
		store := new(fakeLinksStore)

		emptyCode := ""

		params := CreateLinkParams{
			OriginalURL: "https://test.com/test",
			ShortCode:   &emptyCode,
		}

		svc := NewService(store, &fakeShortcodeGen{})

		_, err := svc.CreateLink(t.Context(), params)

		require.ErrorIs(t, err, links.ErrInvlalidShortCode)
	})

	t.Run("create with invalid url", func(t *testing.T) {
		store := new(fakeLinksStore)

		code := "code"

		params := CreateLinkParams{
			OriginalURL: "",
			ShortCode:   &code,
		}

		svc := NewService(store, &fakeShortcodeGen{})

		_, err := svc.CreateLink(t.Context(), params)

		require.ErrorIs(t, err, links.ErrInvlalidURL)
	})

	t.Run("failed to generate unique shortcode", func(t *testing.T) {
		store := new(fakeLinksStore)

		params := CreateLinkParams{
			OriginalURL: "https://test.com/test",
		}

		store.
			On("CreateLink", mock.Anything, mock.Anything).
			Return(links.Link{}, ErrShortCodeConflict)

		generator := &fakeShortcodeGen{codes: []string{"code"}}
		svc := NewService(store, generator)

		_, err := svc.CreateLink(t.Context(), params)

		require.ErrorIs(t, err, ErrFailedToGenerateUniqueShortCode)
		require.Greater(t, generator.calls, 1)
	})
}
