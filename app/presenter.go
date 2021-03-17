package app

import (
	"github.com/chuckha/services/auth/usecases"
)

type UseCasesPresenter struct{}

func (p *UseCasesPresenter) Login(output *usecases.DoLoginOutput) *LoginOutput {
	return &LoginOutput{SessionID: output.SessionID}
}

func (p *UseCasesPresenter) SendLoginMessage(_ *usecases.SendLoginMessageOutput) *SendLoginMessageOutput {
	return &SendLoginMessageOutput{}
}
