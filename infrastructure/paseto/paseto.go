package paseto

import (
	"github.com/o1egl/paseto"
	"github.com/pkg/errors"

	"github.com/chuckha/auth/domain"
)

const (
	oneTimeTokenClaimField = "ott.id"
)

func NewPASETOEncDec() *PASETOEncDec {
	return &PASETOEncDec{paseto.NewV2()}
}

type PASETOEncDec struct {
	encdec *paseto.V2
}

func (p *PASETOEncDec) decrypt(token string, key []byte, loginToken *domain.LoginToken, footer interface{}) error {
	jt := &paseto.JSONToken{}
	if err := p.encdec.Decrypt(token, key, jt, footer); err != nil {
		return errors.WithStack(err)
	}
	loginToken.OneTimeToken.UserID = jt.Subject
	loginToken.NotBefore = jt.NotBefore
	loginToken.Expiration = jt.Expiration
	loginToken.OneTimeToken.Token = jt.Get(oneTimeTokenClaimField)
	return nil
}

func (p *PASETOEncDec) Encode(secretKey string, token *domain.LoginToken, footer interface{}) (string, error) {
	jt := paseto.JSONToken{
		Subject:    token.OneTimeToken.UserID,
		Expiration: token.Expiration,
		NotBefore:  token.NotBefore,
	}
	jt.Set(oneTimeTokenClaimField, token.OneTimeToken.Token)
	return p.encdec.Encrypt([]byte(secretKey), jt, footer)
}

// Decode is the inverse of Encode,
func (p *PASETOEncDec) Decode(secretKey string, token string) (*domain.LoginToken, error) {
	loginToken := &domain.LoginToken{}
	if err := p.decrypt(token, []byte(secretKey), loginToken, nil); err != nil {
		return nil, err
	}
	return loginToken, nil
}
