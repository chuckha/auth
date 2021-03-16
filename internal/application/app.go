package application

import (
	"github.com/chuckha/services/auth/internal/usecases"
)

type Authentication struct {
	UseCases
	Adapter
	Presenter
}

type UseCases interface {
	Login(in *usecases.DoLoginInput) (*usecases.DoLoginOutput, error)
	SendLoginMessage(in *usecases.SendLoginMessageInput) (*usecases.SendLoginMessageOutput, error)
}

type Adapter interface {
	Login(input *LoginInput) *usecases.DoLoginInput
	SendLoginMessage(in *SendLoginMessageInput) *usecases.SendLoginMessageInput
}

type Presenter interface {
	Login(output *usecases.DoLoginOutput) *LoginOutput
	SendLoginMessage(in *usecases.SendLoginMessageOutput) *SendLoginMessageOutput
}

func NewAuthentication(useCases UseCases, adapter Adapter, presenter Presenter) *Authentication {
	return &Authentication{
		UseCases:  useCases,
		Adapter:   adapter,
		Presenter: presenter,
	}
}

type LoginInput struct{}
type LoginOutput struct{}

func (a *Authentication) Login(input *LoginInput) (*LoginOutput, error) {
	out, err := a.UseCases.Login(a.Adapter.Login(input))
	if err != nil {
		return nil, err
	}
	return a.Presenter.Login(out), nil
}

type SendLoginMessageInput struct{}
type SendLoginMessageOutput struct{}

func (a *Authentication) SendLoginMessage(input *SendLoginMessageInput) (*SendLoginMessageOutput, error) {
	out, err := a.UseCases.SendLoginMessage(a.Adapter.SendLoginMessage(input))
	if err != nil {
		return nil, err
	}
	return a.Presenter.SendLoginMessage(out), nil
}
