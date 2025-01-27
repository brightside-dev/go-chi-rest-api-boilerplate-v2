package cmd

import (
	"os"

	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/email"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
)

var testEmailCmd = &cobra.Command{
	Use:   "test_email",
	Short: "Send a test email",
	Run: func(cmd *cobra.Command, args []string) {
		emailService := email.NewEmailService(
			"development",
			os.Getenv("FROM_EMAIL"),
			os.Getenv("FROM_EMAIL_PASSWORD"),
			os.Getenv("FROM_EMAIL_SMTP"),
			os.Getenv("EMAIL_SMTP_ADDRESS"),
		)

		data := &map[string]string{
			"name": "John Doe",
		}
		emailService.SendEmail("test_email", "This is from command", []string{"a@me.com"}, *data)
	},
}

func init() {
	rootCmd.AddCommand(testEmailCmd)
}
