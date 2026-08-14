package functions

import (
	"crypto/rand"
	"math/big"
	"strings"
	"url_shorter_gin/apperrors"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length   = 8
)

func GenerateShortCode() (string, error) {
	var sb strings.Builder
	alphabetLen := big.NewInt(int64(len(alphabet)))

	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", apperrors.ErrGenerateShortCode
		}
		sb.WriteByte(alphabet[index.Int64()])
	}
	return sb.String(), nil
}
