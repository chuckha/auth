package app

import (
	"github.com/chuckha/auth/usecases"
)

type Controller struct {
	Adapter
	UseCases
	Presenter
}

type Adapter interface {
	Login(input *LoginInput) *usecases.DoLoginInput
	SendLoginMessage(in *SendLoginMessageInput) *usecases.SendLoginMessageInput
	CheckValidSession(in *CheckValidSessionInput) *usecases.ValidSessionInput
}

type UseCases interface {
	Login(in *usecases.DoLoginInput) (*usecases.DoLoginOutput, error)
	SendLoginMessage(in *usecases.SendLoginMessageInput) (*usecases.SendLoginMessageOutput, error)
	CheckValidSession(in *usecases.ValidSessionInput) (*usecases.ValidSessionOutput, error)
}

type Presenter interface {
	Login(output *usecases.DoLoginOutput) *LoginOutput
	SendLoginMessage(in *usecases.SendLoginMessageOutput) *SendLoginMessageOutput
	CheckValidSession(output *usecases.ValidSessionOutput) *CheckValidSessionOutput
}

func NewController(useCases UseCases, adapter Adapter, presenter Presenter) *Controller {
	return &Controller{
		UseCases:  useCases,
		Adapter:   adapter,
		Presenter: presenter,
	}
}

type LoginInput struct {
	EncodedToken string
}
type LoginOutput struct {
	SessionID string
}

func (c *Controller) Login(input *LoginInput) (*LoginOutput, error) {
	out, err := c.UseCases.Login(c.Adapter.Login(input))
	if err != nil {
		return nil, err
	}
	return c.Presenter.Login(out), nil
}

type SendLoginMessageInput struct {
	LoginMessageDestination string
}
type SendLoginMessageOutput struct{}

func (c *Controller) SendLoginMessage(input *SendLoginMessageInput) (*SendLoginMessageOutput, error) {
	out, err := c.UseCases.SendLoginMessage(c.Adapter.SendLoginMessage(input))
	if err != nil {
		return nil, err
	}
	return c.Presenter.SendLoginMessage(out), nil
}

type CheckValidSessionInput struct {
	SessionID string
}
type CheckValidSessionOutput struct {
	UID string
}

func (c *Controller) CheckValidSession(in *CheckValidSessionInput) (*CheckValidSessionOutput, error) {
	out, err := c.UseCases.CheckValidSession(c.Adapter.CheckValidSession(in))
	if err != nil {
		return nil, err
	}
	return c.Presenter.CheckValidSession(out), nil
}
