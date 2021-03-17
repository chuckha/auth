package app

import (
	"github.com/chuckha/auth/usecases"
)

type AuthUseCases struct {
	*usecases.DoLogin
	*usecases.LoginMessageSender
}
