package usecases

import (
	"github.com/chuckha/services/auth/domain"
	"github.com/chuckha/services/auth/usecases/dto"
)

type Encoder interface {
	Encode(secretKey string, token *dto.LoginToken, footer interface{}) (string, error)
}

type MessageSender interface {
	SendMessage(destination, contents string) error
}

type MessageGenerator interface {
	GenerateLoginMessage(token string) string
}

type SecretKeyGetter interface {
	GetSecretKey() string
}

type TokenRepository interface {
	GetToken(uid, token string) (*dto.OneTimeToken, error)
	SaveToken(token *dto.OneTimeToken) error
	DeleteToken(uid, token string) error
}

type UserRepository interface {
	GetUser(uid string) (*dto.User, error)
	CreateUser(*dto.User) error
}

type IDGenerator interface {
	ID() string
}

type LoginMessageSender struct {
	IDGenerator
	SecretKeyGetter
	Encoder
	MessageGenerator
	MessageSender
	TokenRepository
	UserRepository
}

type SendLoginMessageInput struct {
	Destination string
}

type SendLoginMessageOutput struct{}

func (l *LoginMessageSender) SendLoginMessage(in *SendLoginMessageInput) (*SendLoginMessageOutput, error) {
	// create a one-time-token and save it
	token := l.ID()
	ott, err := domain.CreateOneTimeToken(in.Destination, token)
	if err != nil {
		return nil, err
	}
	dtoOTT := &dto.OneTimeToken{
		UserID:  ott.UserID,
		Token:   ott.Token,
		Expires: ott.Expires,
	}
	if err := l.SaveToken(dtoOTT); err != nil {
		return nil, err
	}

	// get or create the user
	if _, err = l.GetUser(in.Destination); err != nil {
		newUser, err2 := domain.NewUser(in.Destination)
		if err2 != nil {
			return nil, err
		}
		u := &dto.User{ID: newUser.GetDestination()}
		if err := l.CreateUser(u); err != nil {
			return nil, err
		}
	}

	// create a new LOGIN token using the one time token
	loginToken, err := domain.CreateLoginToken(ott)
	if err != nil {
		return nil, err
	}
	dtoLoginToken := &dto.LoginToken{
		OneTimeToken: ott.Token,
		UserID:       ott.UserID,
		Expiration:   loginToken.GetExpiration(),
		NotBefore:    loginToken.GetNotBefore(),
	}
	encodedToken, err := l.Encode(l.GetSecretKey(), dtoLoginToken, nil)
	if err != nil {
		return nil, err
	}
	contents := l.GenerateLoginMessage(encodedToken)
	err = l.SendMessage(in.Destination, contents)
	return &SendLoginMessageOutput{}, err
}
