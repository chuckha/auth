package domain

import (
	"github.com/pkg/errors"
)

// TODO: configuration
const secretKeySize = 32

type secretKey string

func NewSecretKey(key string) (secretKey, error) {
	if len(key) != secretKeySize {
		return "", errors.Errorf("secret key is not the correct key size of %d", secretKeySize)
	}
	return secretKey(key), nil
}
