package messages

import (
	"fmt"
)

type Generator struct{}

func (g *Generator) GenerateLoginMessage(token string) string {
	return fmt.Sprintf("click this link to login http://localhost:8888/validate_link?token=%s", token)
}
