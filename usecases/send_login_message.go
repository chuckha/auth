package usecases

import (
	"github.com/chuckha/auth/domain"
)

type Encoder interface {
	Encode(secretKey string, token *domain.LoginToken, footer interface{}) (string, error)
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
	GetToken(uid, token string) (*domain.OneTimeToken, error)
	SaveToken(token *domain.OneTimeToken) error
	DeleteToken(uid, token string) error
}

type UserRepository interface {
	GetUser(uid string) (*domain.User, error)
	CreateUser(*domain.User) error
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
	tkn, err := domain.NewOneTimeToken(ott.UserID, ott.Token, ott.Expires)
	if err != nil {
		return nil, err
	}
	if err := l.SaveToken(tkn); err != nil {
		return nil, err
	}

	// get or create the user
	if _, err = l.GetUser(in.Destination); err != nil {
		newUser, err2 := domain.NewUser(in.Destination)
		if err2 != nil {
			return nil, err
		}
		if err := l.CreateUser(newUser); err != nil {
			return nil, err
		}
	}

	// create a new LOGIN token using the one time token
	loginToken, err := domain.CreateLoginToken(ott)
	if err != nil {
		return nil, err
	}
	encodedToken, err := l.Encode(l.GetSecretKey(), loginToken, nil)
	if err != nil {
		return nil, err
	}
	contents := l.GenerateLoginMessage(encodedToken)
	err = l.SendMessage(in.Destination, contents)
	return &SendLoginMessageOutput{}, err
}
