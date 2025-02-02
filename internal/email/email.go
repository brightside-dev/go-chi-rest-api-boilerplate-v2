package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"net/smtp"
	"os"

	Client "github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/email/clients"

	_ "github.com/joho/godotenv/autoload"
)

type EmailServiceInterface interface {
	Send(to []string, subject string, body string) error
	smptSend(to []string, subject string, htmlBody string) error
}

type EmailService struct {
	Env       string
	Logger    *slog.Logger
	EmailAuth *EmailAuth
	Mailgun   *Client.Mailgun
}

type EmailAuth struct {
	FromEmail         string
	FromEmailPassword string
	SMTPHost          string
	SMTPAddr          string
}

func NewEmailService(
	logger *slog.Logger,
) *EmailService {
	return &EmailService{
		Env:    os.Getenv("APP_ENV"),
		Logger: logger,
		EmailAuth: &EmailAuth{
			FromEmail:         os.Getenv("FROM_EMAIL"),
			FromEmailPassword: os.Getenv("FROM_EMAIL_PASSWORD"),
			SMTPHost:          os.Getenv("SMTP_HOST"),
			SMTPAddr:          os.Getenv("SMTP_ADDR"),
		},
		Mailgun: Client.NewMailgun(),
	}
}

func (s *EmailService) Send(
	templateName string,
	subject string,
	to []string,
	data map[string]string) error {

	// Parse and render the HTML template
	tmpl, err := template.ParseFiles("internal/email/templates/" + templateName + ".html")
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	// Conditionally send email based on environment
	if s.Env == "local" {
		// Send email using MailCatcher
		if err := s.localSend(to, subject, rendered.String()); err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
	} else {
		// Send email using Mailgun
		if _, err := s.Mailgun.SendEmail(s.EmailAuth.FromEmail, subject, to[0], rendered); err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
	}

	// Log

	return nil
}

func (s *EmailService) localSend(to []string, subject string, htmlBody string) error {
	var auth smtp.Auth = nil
	if s.Env != "local" {
		auth = smtp.PlainAuth("", s.EmailAuth.FromEmail, "", s.EmailAuth.SMTPHost)
	}

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	message := "Subject: " + subject + "\n" + headers + "\n\n" + htmlBody
	err := smtp.SendMail(
		s.EmailAuth.SMTPAddr,
		auth,
		s.EmailAuth.FromEmail,
		to,
		[]byte(message),
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *EmailService) log(to []string, subject string, htmlBody string) {
	s.Logger.Info("Email sent", "to", to, "subject", subject, "body", htmlBody)
}
