package dto

import (
	"time"
)

type Session struct {
	ID      string
	UserID  string
	Expires time.Time
}

type LoginToken struct {
	OneTimeToken string
	UserID       string
	Expiration   time.Time
	NotBefore    time.Time
}

type OneTimeToken struct {
	UserID  string
	Token   string
	Expires time.Time
}

type User struct {
	ID string
}
