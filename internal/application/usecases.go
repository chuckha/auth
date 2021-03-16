package application

import (
	"github.com/chuckha/services/auth/internal/usecases"
)

type AuthUseCases struct {
	*usecases.DoLogin
	*usecases.LoginMessageSender
}
