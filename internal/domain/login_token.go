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

func (l *loginToken) GetExpiration() time.Time {
	return l.expiration
}

func (l *loginToken) GetNotBefore() time.Time {
	return l.expiration
}

func (l *loginToken) GetOneTimeToken() oneTimeToken {
	return l.oneTimeToken
}

func NewLoginToken(token *oneTimeToken, expires, notBefore time.Time) (*loginToken, error) {
	if time.Now().After(expires) {
		return nil, errors.New("login token is expired")
	}
	if time.Now().Before(notBefore) {
		return nil, errors.New("login token is not yet valid")
	}
	return &loginToken{
		oneTimeToken: *token,
		expiration:   expires,
		notBefore:    notBefore,
	}, nil
}

func CreateLoginToken(token *oneTimeToken) (*loginToken, error) {
	return NewLoginToken(token, time.Now().Add(tokenLifespan), time.Now())
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
