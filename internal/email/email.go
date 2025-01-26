package email

import (
	"fmt"
	"net/smtp"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func SendEmail(to []string, subject string, body string) error {
	fromEmail := os.Getenv("FROM_EMAIL")
	smtpHost := os.Getenv("FROM_EMAIL_SMTP")
	smtpAddr := os.Getenv("EMAIL_SMTP_ADDRESS")

	// Debugging statements
	fmt.Println("FROM_EMAIL:", fromEmail)
	fmt.Println("FROM_EMAIL_SMTP:", smtpHost)
	fmt.Println("EMAIL_SMTP_ADDRESS:", smtpAddr)

	if fromEmail == "" || smtpHost == "" || smtpAddr == "" {
		return fmt.Errorf("missing required environment variables")
	}

	// MailCatcher typically does not require authentication
	//auth := smtp.PlainAuth("", fromEmail, "", smtpHost)

	message := "Subject: " + subject + "\n" + body
	err := smtp.SendMail(
		smtpAddr,
		nil,
		fromEmail,
		to,
		[]byte(message),
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
