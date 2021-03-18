package domain

import (
	"time"

	"github.com/pkg/errors"
)

const (
	//TODO: configuration
	tokenLifespan = 24 * time.Hour
)

type LoginToken struct {
	OneTimeToken OneTimeToken
	Expiration   time.Time
	NotBefore    time.Time
}

func (l *LoginToken) GetExpiration() time.Time {
	return l.Expiration
}

func (l *LoginToken) GetNotBefore() time.Time {
	return l.NotBefore
}

func (l *LoginToken) GetOneTimeToken() OneTimeToken {
	return l.OneTimeToken
}

func NewLoginToken(token *OneTimeToken, expires, notBefore time.Time) (*LoginToken, error) {
	if time.Now().After(expires) {
		return nil, errors.New("login token is expired")
	}
	if time.Now().Before(notBefore) {
		return nil, errors.New("login token is not yet valid")
	}
	if token == nil {
		return nil, errors.New("a login token must have a one time token")
	}
	return &LoginToken{
		OneTimeToken: *token,
		Expiration:   expires,
		NotBefore:    notBefore,
	}, nil
}

func CreateLoginToken(token *OneTimeToken) (*LoginToken, error) {
	return NewLoginToken(token, time.Now().Add(tokenLifespan), time.Now())
}
