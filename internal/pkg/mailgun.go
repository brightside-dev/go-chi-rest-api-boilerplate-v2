package pkg

import (
	"bytes"
	"context"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mailgun/mailgun-go/v4"
)

type Mailgun struct {
	Client *mailgun.MailgunImpl
}

func NewMailgun() *Mailgun {
	domain := os.Getenv("MAILGUN_DOMAIN")
	apiKey := os.Getenv("MAILGUN_API_KEY")

	mg := mailgun.NewMailgun(domain, apiKey)

	return &Mailgun{
		Client: mg,
	}
}

func (m *Mailgun) SendEmail(sender, subject, recipient string, template bytes.Buffer) (string, error) {

	body := template.String()

	message := mailgun.NewMessage(sender, subject, body, recipient)
	message.SetHTML(body)

	ctx := context.Background()
	_, id, err := m.Client.Send(ctx, message)
	if err != nil {
		return "", err
	}

	return id, err
}
