package usecases

import (
	"github.com/chuckha/auth/domain"
)

type ValidSessionInput struct {
	SID string
}
type ValidSessionOutput struct {
	UID string
	SID string
}
type SessionChecker struct {
	LookupSessionRepository
}

func (c *SessionChecker) CheckValidSession(in *ValidSessionInput) (*ValidSessionOutput, error) {
	session, err := c.LookupSession(in.SID)
	if err != nil {
		return nil, err
	}
	domainSesh, err := domain.NewSession(session.ID, session.UserID, session.Expires)
	if err != nil {
		return nil, err
	}
	return &ValidSessionOutput{
		UID: domainSesh.GetUID(),
		SID: domainSesh.GetID(),
	}, nil
}
