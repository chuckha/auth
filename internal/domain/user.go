package domain

import "errors"

type user struct {
	// destination is a unique string to contact a user. could be an email, a phone number or something that uniquely
	// identifies this user. TODO: phone numbers are not necessarily good here because phone numbers can be recycled.
	destination string
}

func NewUser(destination string) (*user, error) {
	if destination == "" {
		return nil, errors.New("user must have an id")
	}
	return &user{destination}, nil
}

func (u *user) GetDestination() string {
	return u.destination
}
