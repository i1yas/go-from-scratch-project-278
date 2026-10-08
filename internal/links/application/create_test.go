package application_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
	"hexleturlshort/internal/links/application/testutils"
)

func TestCreateLink(t *testing.T) {
	t.Run("create with valid shortcode", func(t *testing.T) {
		store := new(fakeLinksStore)

		codeRaw := "test"

		want := testutils.Link(t, codeRaw)

		params := application.CreateLinkParams{
			OriginalURL: string(want.OriginalURL),
			ShortCode:   &codeRaw,
		}

		store.
			On("CreateLink", mock.Anything, want).
			Return(want, nil)

		generator := &fakeShortcodeGen{}
		svc := application.NewService(
			store,
			&fakeVisitsStore{},
			generator,
		)

		got, err := svc.CreateLink(t.Context(), params)

		require.NoError(t, err)
		require.Equal(t, want, got)
		require.Equal(t, 1, len(store.Calls))
		require.Equal(t, 0, generator.calls)
	})

	t.Run("create without shortcode", func(t *testing.T) {
		store := new(fakeLinksStore)

		generatedCode := "generated"

		want := testutils.Link(t, generatedCode)

		params := application.CreateLinkParams{
			OriginalURL: string(want.OriginalURL),
		}

		store.
			On("CreateLink", mock.Anything, want).
			Return(want, nil)

		generator := &fakeShortcodeGen{codes: []string{generatedCode}}
		svc := application.NewService(
			store,
			&fakeVisitsStore{},
			generator,
		)

		got, err := svc.CreateLink(t.Context(), params)

		require.NoError(t, err)
		require.Equal(t, want, got)
		require.Equal(t, 1, len(store.Calls))
		require.Equal(t, 1, generator.calls)
	})

	t.Run("create without shortcode with retry", func(t *testing.T) {
		store := new(fakeLinksStore)

		codes := []string{"code-one", "code-two"}

		want := testutils.Link(t, codes[1])

		params := application.CreateLinkParams{
			OriginalURL: string(want.OriginalURL),
		}

		store.
			On("CreateLink", mock.Anything, links.Link{
				OriginalURL: want.OriginalURL,
				ShortCode:   links.ShortCode(codes[0]),
			}).
			Return(links.Link{}, application.ErrShortCodeConflict).
			Once().
			On("CreateLink", mock.Anything, links.Link{
				OriginalURL: want.OriginalURL,
				ShortCode:   links.ShortCode(codes[1]),
			}).
			Return(want, nil)

		generator := &fakeShortcodeGen{codes: codes}
		svc := application.NewService(
			store,
			&fakeVisitsStore{},
			generator,
		)

		got, err := svc.CreateLink(t.Context(), params)

		require.NoError(t, err)
		require.Equal(t, want, got)
		require.Equal(t, 2, len(store.Calls))
		require.Equal(t, 2, generator.calls)
	})

	t.Run("create with invalid shortcode", func(t *testing.T) {
		store := new(fakeLinksStore)

		emptyCode := ""

		params := application.CreateLinkParams{
			OriginalURL: "https://test.com/test",
			ShortCode:   &emptyCode,
		}

		svc := application.NewService(
			store,
			&fakeVisitsStore{},
			&fakeShortcodeGen{},
		)

		_, err := svc.CreateLink(t.Context(), params)

		require.ErrorContains(t, err, "shortcode:")
		require.ErrorIs(t, err, links.ErrInvlalidShortCode)
		require.Equal(t, 0, len(store.Calls))
	})

	t.Run("create with invalid url", func(t *testing.T) {
		store := new(fakeLinksStore)

		code := "code"

		params := application.CreateLinkParams{
			OriginalURL: "",
			ShortCode:   &code,
		}

		svc := application.NewService(
			store,
			&fakeVisitsStore{},
			&fakeShortcodeGen{},
		)

		_, err := svc.CreateLink(t.Context(), params)

		require.ErrorContains(t, err, "original_url:")
		require.ErrorIs(t, err, links.ErrInvlalidURL)
		require.Equal(t, 0, len(store.Calls))
	})

	t.Run("failed to generate shortcode", func(t *testing.T) {
		store := new(fakeLinksStore)

		params := application.CreateLinkParams{
			OriginalURL: "https://test.com/test",
		}

		generator := &fakeShortcodeGen{}
		svc := application.NewService(
			store,
			&fakeVisitsStore{},
			generator,
		)

		_, err := svc.CreateLink(t.Context(), params)

		require.ErrorIs(t, err, application.ErrShortCodeGeneratorInternal)
		require.Equal(t, 0, len(store.Calls))
		require.Equal(t, 1, generator.calls)
	})

	t.Run("failed to generate unique shortcode", func(t *testing.T) {
		store := new(fakeLinksStore)

		params := application.CreateLinkParams{
			OriginalURL: "https://test.com/test",
		}

		store.
			On("CreateLink", mock.Anything, mock.Anything).
			Return(links.Link{}, application.ErrShortCodeConflict)

		generator := &fakeShortcodeGen{codes: []string{"code"}}
		svc := application.NewService(
			store,
			&fakeVisitsStore{},
			generator,
		)

		_, err := svc.CreateLink(t.Context(), params)

		require.ErrorIs(t, err, application.ErrFailedToGenerateUniqueShortCode)
		require.Greater(t, len(store.Calls), 2)
		require.Greater(t, generator.calls, 2)
		require.Equal(t, len(store.Calls), generator.calls)
	})

	t.Run("generated shortcode is invalid", func(t *testing.T) {
		store := new(fakeLinksStore)

		params := application.CreateLinkParams{
			OriginalURL: "https://test.com/test",
		}

		generator := &fakeShortcodeGen{codes: []string{""}}
		svc := application.NewService(
			store,
			&fakeVisitsStore{},
			generator,
		)

		_, err := svc.CreateLink(t.Context(), params)

		require.ErrorIs(t, err, application.ErrShortCodeGeneratorInternal)
		require.Equal(t, generator.calls, 1)
	})
}
