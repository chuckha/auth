package app

import (
	"github.com/chuckha/services/auth/usecases"
)

type AuthUseCases struct {
	*usecases.DoLogin
	*usecases.LoginMessageSender
}
