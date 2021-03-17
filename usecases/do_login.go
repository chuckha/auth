package usecases

import (
	"github.com/pkg/errors"

	"github.com/chuckha/services/auth/domain"
	"github.com/chuckha/services/auth/usecases/dto"
)

/*
// Encode returns a signed & encrypted single-use token with claims.


type Encryptor interface {
	Encrypt([]byte, interface{}, interface{}) (string, error)
}

type SetClaimer interface {
	Set(string, string)
}


*/
type SetClaimer interface {
	Set(string, string)
}

type Decoder interface {
	Decode(secretKey string, token string) (*dto.LoginToken, error)
}

type SessionRepository interface {
	GetSession(namespace, id string) (*dto.Session, error)
	SaveSession(session *dto.Session) error
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
	dtoToken, err := d.Decode(sk, in.EncodedLoginToken)
	if err != nil {
		return nil, err
	}
	// regardless, we've used it. its one use is over.
	defer d.DeleteToken(dtoToken.UserID, dtoToken.OneTimeToken)

	// ensure the token has not yet been used
	foundToken, err := d.GetToken(dtoToken.UserID, dtoToken.OneTimeToken)
	if err != nil {
		return nil, err
	}

	// ensure the one time token has valid data
	oneTimeToken, err := domain.NewOneTimeToken(foundToken.UserID, foundToken.Token, foundToken.Expires)
	if err != nil {
		return nil, err
	}

	// ensure the login token was also valid
	if _, err := domain.NewLoginToken(oneTimeToken, dtoToken.Expiration, dtoToken.NotBefore); err != nil {
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
