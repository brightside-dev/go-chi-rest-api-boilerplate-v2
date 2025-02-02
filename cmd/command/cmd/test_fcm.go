package cmd

import (
	"fmt"
	"log"

	Client "github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/push/clients"

	"firebase.google.com/go/messaging"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
)

var testFcmCmd = &cobra.Command{
	Use:   "test_fcm",
	Short: "Send a test email",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize FCM client
		fcm, err := Client.NewFCM()
		if err != nil {
			log.Fatalf("Failed to initialize FCM: %v", err)
		}

		// Replace with an actual FCM device token
		testToken := "your_device_fcm_token"

		// Create a test message
		message := &messaging.Message{
			Token: testToken,
			Notification: &messaging.Notification{
				Title: "Test Notification",
				Body:  "This is a test message from Go!",
			},
		}

		// Send the message
		err = fcm.Push(message)
		if err != nil {
			log.Fatalf("Failed to send push notification: %v", err)
		}

		fmt.Println("Push notification sent successfully!")
	},
}

func init() {
	rootCmd.AddCommand(testFcmCmd)
}
