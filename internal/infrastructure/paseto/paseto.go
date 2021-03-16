package paseto

import (
	"time"

	"github.com/o1egl/paseto"
	"github.com/pkg/errors"

	"github.com/chuckha/services/auth/internal/usecases/dto"
)

type PASETO struct {
	paseto.JSONToken
}

func NewSetClaimer(token *dto.LoginToken) *PASETO {
	return &PASETO{
		paseto.JSONToken{
			Subject:    token.UserID,
			IssuedAt:   time.Now(),
			Expiration: token.Expiration,
			NotBefore:  token.NotBefore,
		},
	}

}

func (p *PASETO) SetClaim(key, value string) {
	p.JSONToken.Set(key, value)
}

func NewPASETOEncDec() *PASETOEncDec {
	return &PASETOEncDec{paseto.NewV2()}
}

type PASETOEncDec struct {
	encdec *paseto.V2
}

func (p *PASETOEncDec) Decrypt(token string, key []byte, loginToken *dto.LoginToken, footer interface{}) error {
	jt := &paseto.JSONToken{}
	if err := p.encdec.Decrypt(token, key, jt, footer); err != nil {
		return errors.WithStack(err)
	}
	loginToken.UserID = jt.Get("uid")
	loginToken.NotBefore = jt.NotBefore
	loginToken.Expiration = jt.Expiration
	loginToken.OneTimeToken = jt.Get("ott.id")
	return nil
}

func (p *PASETOEncDec) Encrypt(key []byte, payload interface{}, footer interface{}) (string, error) {
	return p.encdec.Encrypt(key, payload, footer)
}
