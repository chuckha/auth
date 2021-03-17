package domain

import (
	"time"

	"github.com/pkg/errors"
)

const (
	// TODO: configuration
	sessionDuration = 30 * 24 * time.Hour
)

type session struct {
	id      string
	uid     string
	expires time.Time
}

func CreateSession(sid, uid string) (*session, error) {
	return NewSession(sid, uid, time.Now().Add(sessionDuration))
}

func NewSession(sid, uid string, expires time.Time) (*session, error) {
	if sid == "" {
		return nil, errors.New("session ID cannot be empty")
	}
	if uid == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	if time.Now().After(expires) {
		return nil, errors.New("session is expired")
	}
	return &session{
		id:      sid,
		uid:     uid,
		expires: expires,
	}, nil
}

func (s *session) GetExpires() time.Time {
	return s.expires
}

func (s *session) GetID() string {
	return s.id
}

func (s *session) GetUID() string {
	return s.uid
}
