package cmd

import (
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
)

var testMailgunCmd = &cobra.Command{
	Use:   "test_mailgun",
	Short: "Send a test email",
	Run: func(cmd *cobra.Command, args []string) {
		if container == nil {
			log.Fatal("Container is not initialized")
		}
		data := &map[string]string{
			"name": "John Doe",
		}

		err := container.EmailService.SendEmail("test_email", "This is from command", []string{"battousai.dev@proton.me"}, *data)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(testMailgunCmd)
}
