package domain

import (
	"time"

	"github.com/pkg/errors"
)

const (
	//TODO: configuration
	tokenLifespan          = 24 * time.Hour
	oneTimeTokenIDClaimKey = "ott.id"
	uidClaimKey            = "uid"
)

type loginToken struct {
	oneTimeToken oneTimeToken
	expiration   time.Time
	notBefore    time.Time
	uid          string
}

func (l *loginToken) GetExpiration() string {
	return l.expiration.Format(time.RFC3339)
}

func (l *loginToken) GetNotBefore() string {
	return l.expiration.Format(time.RFC3339)
}

func (l *loginToken) GetOneTimeToken() oneTimeToken {
	return l.oneTimeToken
}

func NewLoginToken(token *oneTimeToken, expiration, notBefore string) (*loginToken, error) {
	expires, err := time.Parse(time.RFC3339, expiration)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if time.Now().After(expires) {
		return nil, errors.New("login token is expired")
	}
	validAfter, err := time.Parse(time.RFC3339, notBefore)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if time.Now().Before(validAfter) {
		return nil, errors.New("login token is not yet valid")
	}
	return &loginToken{
		oneTimeToken: *token,
		expiration:   expires,
		notBefore:    validAfter,
	}, nil
}

func CreateLoginToken(token *oneTimeToken) (*loginToken, error) {
	return NewLoginToken(token, time.Now().Add(tokenLifespan).Format(time.RFC3339), time.Now().Format(time.RFC3339))
}

// Encode returns a signed & encrypted single-use token with claims.
func (lt loginToken) Encode(secretKey secretKey, encryptor Encryptor, claim SetClaimer) (string, error) {
	claim.Set(oneTimeTokenIDClaimKey, lt.oneTimeToken.Token)
	claim.Set(uidClaimKey, lt.uid)
	return encryptor.Encrypt([]byte(secretKey), claim, nil)
}

// Decode is the inverse of Encode,
func Decode(secretKey secretKey, decryptor Decryptor, token string) (*loginToken, error) {
	loginToken := &loginToken{}
	if err := decryptor.Decrypt(token, []byte(secretKey), loginToken, nil); err != nil {
		return nil, err
	}
	return NewLoginToken(&loginToken.oneTimeToken, loginToken.GetExpiration(), loginToken.GetNotBefore())
}

/* dependencies */

type Encryptor interface {
	Encrypt([]byte, interface{}, interface{}) (string, error)
}

type SetClaimer interface {
	Set(string, string)
}

type Decryptor interface {
	Decrypt(string, []byte, *loginToken, interface{}) error
}
