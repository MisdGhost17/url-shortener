package codegen

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"url-shortener/internal/core/domain"
)

type RandomGenerator struct{}

func (RandomGenerator) Generate() (string, error) {
	code := make([]byte, domain.CodeLength)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(domain.CodeAlphabet))))
		if err != nil {
			return "", fmt.Errorf("generate random code: %w", err)
		}
		code[i] = domain.CodeAlphabet[n.Int64()]
	}
	return string(code), nil
}
