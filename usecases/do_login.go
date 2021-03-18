package usecases

import (
	"github.com/pkg/errors"

	"github.com/chuckha/auth/domain"
)

type Decoder interface {
	Decode(secretKey string, token string) (*domain.LoginToken, error)
}

type SessionRepository interface {
	GetSession(uid, id string) (*domain.Session, error)
	SaveSession(session *domain.Session) error
}

type LookupSessionRepository interface {
	LookupSession(id string) (*domain.Session, error)
}

type DoLogin struct {
	IDGenerator
	SecretKeyGetter
	Decoder
	TokenRepository
	SessionRepository
}

type DoLoginInput struct {
	EncodedLoginToken string
}
type DoLoginOutput struct {
	SessionID string
}

func (d *DoLogin) Login(in *DoLoginInput) (*DoLoginOutput, error) {
	sk := d.GetSecretKey()
	token, err := d.Decode(sk, in.EncodedLoginToken)
	if err != nil {
		return nil, err
	}
	// regardless, we've used it. its one use is over.
	defer d.DeleteToken(token.OneTimeToken.UserID, token.OneTimeToken.Token)

	// ensure the token has not yet been used
	foundToken, err := d.GetToken(token.OneTimeToken.UserID, token.OneTimeToken.Token)
	if err != nil {
		return nil, err
	}

	// ensure the one time token has valid data
	oneTimeToken, err := domain.NewOneTimeToken(foundToken.UserID, foundToken.Token, foundToken.Expires)
	if err != nil {
		return nil, err
	}

	// ensure the login token was also valid
	if _, err := domain.NewLoginToken(oneTimeToken, token.Expiration, token.NotBefore); err != nil {
		return nil, err
	}

	// create login session
	sid := d.ID()

	// make sure the new session is unique
	if _, err := d.GetSession(oneTimeToken.UserID, sid); err == nil {
		return nil, errors.New("session identifier already exists for user, highly unlikely this error ever happens")
	}
	session, err := domain.CreateSession(sid, oneTimeToken.UserID)
	if err != nil {
		return nil, err
	}

	// Save the new session
	if err := d.SaveSession(session); err != nil {
		return nil, err
	}

	// return the session ID
	return &DoLoginOutput{
		SessionID: session.GetID(),
	}, nil

}
