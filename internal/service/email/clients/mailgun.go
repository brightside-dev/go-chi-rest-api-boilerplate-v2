package email_clients

import (
	"bytes"
	"context"
	"os"

	_ "github.com/joho/godotenv/autoload"
	mg "github.com/mailgun/mailgun-go/v4"
)

type Mailgun interface {
	SendEmail(sender, subject, recipient string, template bytes.Buffer) (string, error)
}

type mailgun struct {
	Client *mg.MailgunImpl
}

func NewMailgun() Mailgun {
	domain := os.Getenv("MAILGUN_DOMAIN")
	apiKey := os.Getenv("MAILGUN_API_KEY")

	mg := mg.NewMailgun(domain, apiKey)

	return &mailgun{
		Client: mg,
	}
}

func (m *mailgun) SendEmail(sender, subject, recipient string, template bytes.Buffer) (string, error) {

	body := template.String()

	message := mg.NewMessage(sender, subject, body, recipient)
	message.SetHTML(body)

	ctx := context.Background()
	_, id, err := m.Client.Send(ctx, message)
	if err != nil {
		return "", err
	}

	return id, err
}
