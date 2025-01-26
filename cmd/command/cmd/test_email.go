package cmd

import (
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/email"

	"github.com/spf13/cobra"
)

var testEmailCmd = &cobra.Command{
	Use:   "test_email",
	Short: "Send a test email",
	Run: func(cmd *cobra.Command, args []string) {
		email.SendEmail([]string{"a@me.com"}, "Test email", "This is a test email")
	},
}

func init() {
	rootCmd.AddCommand(testEmailCmd)
}
