package shortcodegen

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"

	"hexleturlshort/internal/links/application"
)

// Generator contains randSource for shortcode generation
type Generator struct {
	randSource io.Reader
}

// DefaultRandSource should be used outside of tests
var DefaultRandSource = rand.Reader

// NewGenerator creates Generator
func NewGenerator(randSource io.Reader) *Generator {
	return &Generator{
		randSource: randSource,
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Generate generates random string for shortcode
func (g *Generator) Generate() (string, error) {
	result := make([]byte, 32)
	maxInd := big.NewInt(int64(len(charset)))

	for i := range result {
		charInd, err := rand.Int(g.randSource, maxInd)
		if err != nil {
			return "", fmt.Errorf("%w: %w",
				application.ErrShortCodeGeneratorInternal,
				err,
			)
		}

		result[i] = charset[charInd.Int64()]
	}

	return string(result), nil
}
