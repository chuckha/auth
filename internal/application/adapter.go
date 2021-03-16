package application

import (
	"github.com/chuckha/services/auth/internal/usecases"
)

type UseCasesAdapter struct{}

func (a *UseCasesAdapter) Login(input *LoginInput) *usecases.DoLoginInput {
	return nil
}

func (a *UseCasesAdapter) SendLoginMessage(in *SendLoginMessageInput) *usecases.SendLoginMessageInput {
	return nil
}
