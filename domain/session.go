package domain

import (
	"time"

	"github.com/pkg/errors"
)

const (
	// TODO: configuration
	sessionDuration = 30 * 24 * time.Hour
)

type Session struct {
	ID      string
	UID     string
	Expires time.Time
}

func CreateSession(sid, uid string) (*Session, error) {
	return NewSession(sid, uid, time.Now().Add(sessionDuration))
}

func NewSession(sid, uid string, expires time.Time) (*Session, error) {
	if sid == "" {
		return nil, errors.New("Session ID cannot be empty")
	}
	if uid == "" {
		return nil, errors.New("User ID cannot be empty")
	}
	if time.Now().After(expires) {
		return nil, errors.New("Session is expired")
	}
	return &Session{
		ID:      sid,
		UID:     uid,
		Expires: expires,
	}, nil
}

func (s *Session) GetExpires() time.Time {
	return s.Expires
}

func (s *Session) GetID() string {
	return s.ID
}

func (s *Session) GetUID() string {
	return s.UID
}
