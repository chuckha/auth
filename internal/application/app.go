package application

import (
	"github.com/chuckha/services/auth/internal/usecases"
)

type Authentication struct {
	UseCases
	Adapter
	Presenter
}

type Adapter interface {
	Login(input *LoginInput) *usecases.DoLoginInput
	SendLoginMessage(in *SendLoginMessageInput) *usecases.SendLoginMessageInput
}

type Presenter interface {
	Login(output *usecases.DoLoginOutput) *LoginOutput
	SendLoginMessage(in *usecases.SendLoginMessageOutput) *SendLoginMessageOutput
}

type UseCases interface {
	Login(in *usecases.DoLoginInput) (*usecases.DoLoginOutput, error)
	SendLoginMessage(in *usecases.SendLoginMessageInput) (*usecases.SendLoginMessageOutput, error)
}

// Think of the input as the controller input and the output is a view model

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
