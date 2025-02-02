package push_client

import (
	"fmt"
	"os"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"

	_ "github.com/joho/godotenv/autoload"
)

type APN struct {
	Client *apns2.Client
}

func NewAPN() (*APN, error) {
	// Load APN certificate
	cert, err := certificate.FromP12File(os.Getenv("APN_CERT_PATH"), "APN_CERT_PASSWORD")
	if err != nil {
		return nil, fmt.Errorf("failed to load APN certificate: %w", err)
	}

	// Create APN client
	client := apns2.NewClient(cert)
	if os.Getenv("APP_ENV") == "production" {
		client = client.Production()
	} else {
		client = client.Development()
	}

	return &APN{Client: client}, nil
}

func (a *APN) Push(notification *apns2.Notification) error {
	res, err := a.Client.Push(notification)
	if err != nil {
		return err
	}

	fmt.Printf("APN Push Success: %v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
	return nil
}

func (a *APN) PushMultiple(notifications []*apns2.Notification) error {
	for _, notification := range notifications {
		if err := a.Push(notification); err != nil {
			return err
		}
	}
	return nil
}
