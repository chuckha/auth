package usecases

import (
	"github.com/pkg/errors"

	"github.com/chuckha/services/auth/internal/domain"
	"github.com/chuckha/services/auth/internal/usecases/dto"
)

type SessionRepository interface {
	GetSession(namespace, id string) (*dto.Session, error)
	SaveSession(session *dto.Session) error
}

type DoLogin struct {
	IDGenerator
	SecretKeyGetter
	domain.Decryptor
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
	sk, err := domain.NewSecretKey(d.GetSecretKey())
	if err != nil {
		return nil, err
	}
	loginToken, err := domain.Decode(sk, d.Decryptor, in.EncodedLoginToken)
	if err != nil {
		return nil, err
	}
	// check if one time token is valid
	ott := loginToken.GetOneTimeToken()
	// if this is successful then the token was good
	oneTimeToken, err := d.GetToken(ott.UserID, ott.Token)
	if err != nil {
		return nil, err
	}

	// remove the token from the repository
	if err := d.DeleteToken(oneTimeToken.UserID, oneTimeToken.Token); err != nil {
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
	dtoSession := &dto.Session{
		ID:      session.GetID(),
		UserID:  session.GetUID(),
		Expires: session.GetExpires(),
	}

	// Save the new session
	if err := d.SaveSession(dtoSession); err != nil {
		return nil, err
	}

	// return the session ID
	return &DoLoginOutput{
		SessionID: session.GetID(),
	}, nil

}
