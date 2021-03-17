package uuid

import (
	"github.com/google/uuid"
)

type UUID struct{}

func (u *UUID) ID() string {
	return uuid.Must(uuid.NewUUID()).String()
}
