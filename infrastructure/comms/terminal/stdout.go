package terminal

import (
	"fmt"
)

type Client struct{}

// Send sends content to the local terminal
func (c *Client) SendMessage(destination, contents string) error {
	fmt.Printf("%s <%s>\n", destination, "auth user")
	fmt.Println("---")
	fmt.Printf("To: %s\n", destination)
	fmt.Printf("Subject: %s\n", "login link inside!")
	fmt.Println(contents)
	fmt.Println("---")
	return nil
}
