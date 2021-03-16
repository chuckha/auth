package dto

type Session struct {
	ID      string
	UserID  string
	Expires string
}

type LoginToken struct {
	OneTimeToken string
	Namespace    string
	Expiration   string
	NotBefore    string
}

type OneTimeToken struct {
	UserID string
	Token  string
}

type User struct {
	Destination string
}
