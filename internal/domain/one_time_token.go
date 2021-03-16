package domain

import (
	"github.com/pkg/errors"
)

type oneTimeToken struct {
	UserID string
	Token  string
}

func NewOneTimeToken(uid, token string) (*oneTimeToken, error) {
	if uid == "" {
		return nil, errors.New("uid cannot be empty")
	}
	if token == "" {
		return nil, errors.New("token cannot be empty")

	}
	return &oneTimeToken{
		UserID: uid,
		Token:  token,
	}, nil
}
