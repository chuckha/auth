package secret

import (
	"strings"
)

type KeyGetter struct {}

func (k *KeyGetter)	GetSecretKey() string {
	return strings.Repeat("8", 32)
}
