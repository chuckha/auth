package app

import (
	"github.com/chuckha/auth/usecases"
)

type UseCasesAdapter struct{}

func (a *UseCasesAdapter) Login(input *LoginInput) *usecases.DoLoginInput {
	return &usecases.DoLoginInput{EncodedLoginToken: input.EncodedToken}
}

func (a *UseCasesAdapter) SendLoginMessage(in *SendLoginMessageInput) *usecases.SendLoginMessageInput {
	return &usecases.SendLoginMessageInput{Destination: in.LoginMessageDestination}
}

func (a *UseCasesAdapter) CheckValidSession(in *CheckValidSessionInput) *usecases.ValidSessionInput {
	return &usecases.ValidSessionInput{
		SID: in.SessionID,
	}
}
