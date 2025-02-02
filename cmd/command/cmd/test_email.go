package cmd

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
)

var testEmailCmd = &cobra.Command{
	Use:   "test_email",
	Short: "Send a test email",
	Run: func(cmd *cobra.Command, args []string) {
		if container == nil {
			log.Fatal("Container is not initialized")
		}

		data := &map[string]string{
			"name": "John Doe",
		}
		err := container.EmailService.Send("test_email", "This is from command", []string{"a@me.com"}, *data)
		if err != nil {
			log.Fatalf("Failed to send email: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(testEmailCmd)
}
