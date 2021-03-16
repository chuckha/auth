package application

import (
	"github.com/chuckha/services/auth/internal/usecases"
)

type UseCasesPresenter struct{}

func (p *UseCasesPresenter) Login(output *usecases.DoLoginOutput) *LoginOutput {
	return nil
}

func (p *UseCasesPresenter) SendLoginMessage(in *usecases.SendLoginMessageOutput) *SendLoginMessageOutput {
	return nil
}
