package domain

import (
	"time"

	"github.com/pkg/errors"
)

const (
	oneTimeTokenLifespan = 5 * time.Hour
)

type oneTimeToken struct {
	UserID  string
	Token   string
	Expires time.Time
}

func CreateOneTimeToken(uid, token string) (*oneTimeToken, error) {
	return NewOneTimeToken(uid, token, time.Now().Add(oneTimeTokenLifespan))
}

func NewOneTimeToken(uid, token string, expires time.Time) (*oneTimeToken, error) {
	if uid == "" {
		return nil, errors.New("uid cannot be empty")
	}
	if token == "" {
		return nil, errors.New("token cannot be empty")
	}
	if time.Now().After(expires) {
		return nil, errors.New("one time token has expired")
	}
	return &oneTimeToken{
		UserID:  uid,
		Token:   token,
		Expires: expires,
	}, nil
}
