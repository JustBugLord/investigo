package investigo

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

var alphabet = "abcdefghijklmnopqrstuvwxyz012345"

type RandomGenerator struct{}

func NewRandomGenerator() *RandomGenerator {
	return &RandomGenerator{}
}

func (rg *RandomGenerator) String(n int) (string, error) {
	if n <= 0 {
		return "", nil
	}
	alphabetLen := len(alphabet)
	randomBytes := make([]byte, n)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	var result strings.Builder
	for i := 0; i < n; i++ {
		idx := int(randomBytes[i]) % alphabetLen
		result.WriteByte(alphabet[idx])
	}
	return result.String(), nil
}

func (rg *RandomGenerator) Number(max int64) (int64, error) {
	if max <= 0 {
		return 0, nil
	}
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return 0, err
	}
	return n.Int64(), nil
}

func (rg *RandomGenerator) NumberString(max int64) (string, error) {
	if max <= 1 {
		return "0", nil
	}
	maxStr := fmt.Sprintf("%d", max-1)
	width := len(maxStr)
	num, err := rg.Number(max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%0*d", width, num), nil
}
