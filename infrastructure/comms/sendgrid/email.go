package sendgrid

import (
	"os"

	"github.com/pkg/errors"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Client is used to send emails via sendgrid.
type Client struct {
	client *sendgrid.Client
	from   *mail.Email
}

type ClientConfig struct {
	APIKey string
	Name string
	From string
	BaseURL string
}

func (c *Client) SendMessage(destination, contents string) error {
	m := mail.NewSingleEmail(c.from, "auth login", emailToSendgridMail(destination), contents, contents)
	resp, err := c.client.Send(m)
	if err != nil {
		return errors.WithStack(err)
	}
	if resp.StatusCode >= 400 {
		// TODO logging in the client
		return errors.Errorf("sendgrid returned %d", resp.StatusCode)
	}
	return nil
}

// NewClient creates a configured client.
// The key can be empty if you set SENDGRID_API_KEY env var to your API key.
func NewClient(config *ClientConfig) (*Client, error) {
	if config.APIKey == "" {
		config.APIKey = os.Getenv("SENDGRID_API_KEY")
	}
	if config.APIKey == "" {
		return nil, errors.New("sendgrid client requires either setting the API in the config or setting the SENDGRID_API_KEY env var")
	}
	sendgridClient := sendgrid.NewSendClient(config.APIKey)
	if config.BaseURL != "" {
		sendgridClient.BaseURL = config.BaseURL
	}

	return &Client{
		client: sendgridClient,
		from: &mail.Email{
			Name:    config.Name,
			Address: config.From,
		},
	}, nil
}

// convert from our types to sendgrid types
func emailToSendgridMail(address string) *mail.Email {
	return &mail.Email{
		Name:    "auth user",
		Address: address,
	}
}
