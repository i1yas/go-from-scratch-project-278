package shortcodegen

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
)

func TestGenerator(t *testing.T) {
	cases := []struct {
		seed string
	}{
		{
			seed: "one",
		},
		{
			seed: "SECRET",
		},
		{
			seed: "abcdefg123",
		},
		{
			seed: "longlonglongseed123",
		},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("seed %s", tc.seed)

		t.Run(name, func(t *testing.T) {
			randSource := testSeededRandSource(tc.seed)
			generator := NewGenerator(randSource)

			got, err := generator.Generate()
			require.NoError(t, err)

			_, err = links.NewShortCode(got)
			require.NoError(t, err)
		})
	}

	t.Run("error broken rand source", func(t *testing.T) {
		brokenRandSource := brokenReader{}
		generator := NewGenerator(&brokenRandSource)

		_, err := generator.Generate()
		require.ErrorIs(t, err, application.ErrFailedToGenerateValidShortCode)
	})
}

func testSeededRandSource(seed string) *rand.ChaCha8 {
	var s [32]byte
	copy(s[:], seed)

	return rand.NewChaCha8(s)
}

type brokenReader struct{}

func (r *brokenReader) Read([]byte) (int, error) {
	return 0, errors.New("broken")
}
