package push_client

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type FCM struct {
	Client *messaging.Client
}

func NewFCM() (*FCM, error) {
	// Initialize Firebase App
	var opt option.ClientOption
	if os.Getenv("APP_ENV") == "production" {
		opt = option.WithCredentialsFile(os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY_PROD"))
	} else {
		opt = option.WithCredentialsFile(os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY_TEST"))
	}

	fmt.Println("Firebase service account key path: ", opt)

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	// Create FCM client
	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize FCM client: %w", err)
	}

	return &FCM{Client: client}, nil
}

func (f *FCM) Push(message *messaging.Message) error {
	// Send notification
	response, err := f.Client.Send(context.Background(), message)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully sent FCM notification: %s\n", response)
	return nil
}

func (f *FCM) PushMultiple(tokens []string, title string, body string) error {
	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	// Send notification
	response, err := f.Client.SendMulticast(context.Background(), message)
	if err != nil {
		return fmt.Errorf("failed to send FCM multicast push: %w", err)
	}

	fmt.Printf("Successfully sent FCM notifications: %d successful, %d failed\n",
		response.SuccessCount, response.FailureCount)

	return nil
}
