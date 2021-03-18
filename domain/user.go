package domain

import "errors"

type User struct {
	// ID is a unique string to contact a User. could be an email, a phone number or something that uniquely
	// identifies this User. TODO: phone numbers are not necessarily good here because phone numbers can be recycled.
	ID string
}

func NewUser(destination string) (*User, error) {
	if destination == "" {
		return nil, errors.New("user must have an identifier")
	}
	return &User{destination}, nil
}

func (u *User) GetDestination() string {
	return u.ID
}
