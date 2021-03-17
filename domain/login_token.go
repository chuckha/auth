package domain

import (
	"time"

	"github.com/pkg/errors"
)

const (
	//TODO: configuration
	tokenLifespan = 24 * time.Hour
)

type loginToken struct {
	oneTimeToken oneTimeToken
	expiration   time.Time
	notBefore    time.Time
}

func (l *loginToken) GetExpiration() time.Time {
	return l.expiration
}

func (l *loginToken) GetNotBefore() time.Time {
	return l.notBefore
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
	if token == nil {
		return nil, errors.New("a login token must have a one time token")
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
