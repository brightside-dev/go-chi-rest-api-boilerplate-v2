package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"

	_ "github.com/joho/godotenv/autoload"
)

type EmailServiceInterface interface {
	SendEmail(to []string, subject string, body string) error
}

type EmailService struct {
	Env       string
	EmailAuth *EmailAuth
}

type EmailAuth struct {
	FromEmail         string
	FromEmailPassword string
	SMTPHost          string
	SMTPAddr          string
}

func NewEmailService(
	env string,
	fromEmail string,
	fromEmailPassword string,
	smtpHost string,
	smtpAddr string,
) *EmailService {
	return &EmailService{
		Env: env,
		EmailAuth: &EmailAuth{
			FromEmail:         fromEmail,
			FromEmailPassword: fromEmailPassword,
			SMTPHost:          smtpHost,
			SMTPAddr:          smtpAddr,
		},
	}
}

func (s *EmailService) SendEmail(templateName string, data map[string]string) error {
	// Parse the HTML template
	tmpl, err := template.ParseFiles("internal/email/templates/" + templateName + ".html")
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	// Render the template with the map data
	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	to, ok := data["to"]
	if !ok {
		return fmt.Errorf("to email address not provided")
	}

	subject, ok := data["subject"]
	if !ok {
		return fmt.Errorf("email subject not provided")
	}

	err = s.smptSend([]string{to}, subject, rendered.String())
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)

	}

	return nil
}

func (s *EmailService) smptSend(to []string, subject string, htmlBody string) error {
	var auth smtp.Auth = nil
	if s.Env != "development" {
		auth = smtp.PlainAuth("", s.EmailAuth.FromEmail, "", s.EmailAuth.SMTPHost)
	}

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	message := "Subject: " + subject + "\n" + headers + "\n\n" + htmlBody
	return smtp.SendMail(
		s.EmailAuth.SMTPAddr,
		auth,
		s.EmailAuth.FromEmail,
		to,
		[]byte(message),
	)
}
