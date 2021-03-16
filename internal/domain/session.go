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
	return NewSession(sid, uid, time.Now().Add(sessionDuration).Format(time.RFC3339))
}

func NewSession(sid, uid string, expiration string) (*session, error) {
	if sid == "" {
		return nil, errors.New("session ID cannot be empty")
	}
	if uid == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	expires, err := time.Parse(time.RFC3339, expiration)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &session{
		id:      sid,
		uid:     uid,
		expires: expires,
	}, nil
}

func (s *session) GetExpires() string {
	return s.expires.Format(time.RFC3339)
}
func (s *session) GetID() string {
	return s.id
}
func (s *session) GetUID() string {
	return s.uid
}
