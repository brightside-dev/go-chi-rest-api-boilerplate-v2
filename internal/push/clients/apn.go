package push_client

import (
	"fmt"
	"log"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

type APN struct {
	Client *apns2.Client
}

func NewAPN() (*APN, error) {
	cert, err := certificate.FromP12File("../cert.p12", "")
	if err != nil {
		log.Fatal("Cert Error:", err)
	}
	client := apns2.NewClient(cert).Production()

	return &APN{Client: client}, nil
}

func (a *APN) Push(notification *apns2.Notification) error {
	res, err := a.Client.Push(notification)

	if err != nil {
		return err
	}

	fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
	return nil
}

func (a *APN) PushMultiple(notifications []*apns2.Notification) error {
	return nil
}
