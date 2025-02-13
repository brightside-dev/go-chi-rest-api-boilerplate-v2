package push

import (
	"log/slog"

	"firebase.google.com/go/messaging"
	Client "github.com/brightside-dev/ronin-fitness-be/internal/service/push/clients"
	"github.com/brightside-dev/ronin-fitness-be/internal/util"
	"github.com/sideshow/apns2"
)

type PushService interface {
	PushIOS()
	PushAndroid(token string, title string, body string)
}

type pushService struct {
	Logger    *slog.Logger
	APNClient Client.APN
	FCMClient Client.FCM
}

func NewPushService(logger *slog.Logger) (PushService, error) {
	apnClient, err := Client.NewAPN()
	if err != nil {
		return nil, err
	}

	fcmClient, err := Client.NewFCM()
	if err != nil {
		return nil, err
	}

	return &pushService{
		Logger:    logger,
		APNClient: apnClient,
		FCMClient: fcmClient,
	}, nil
}

func (p *pushService) PushIOS() {
	notification := &apns2.Notification{
		DeviceToken: "device_token",
		Topic:       "com.example.app",
		Payload:     []byte(`{"aps":{"alert":"Hello!"}}`),
	}

	err := p.APNClient.Push(notification)
	if err != nil {
		context := map[string]interface{}{
			"device_token": &notification.DeviceToken,
			"topic":        &notification.Topic,
			"payload":      &notification.Payload,
			"error":        err.Error(),
		}

		util.LogWithContext(
			p.Logger,
			slog.LevelError,
			"Failed to push notification",
			context,
			nil,
		)
	}
}

func (p *pushService) PushAndroid(token string, title string, body string) {
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	err := p.FCMClient.Push(message)
	if err != nil {
		context := map[string]interface{}{
			"token": token,
			"title": title,
			"body":  body,
			"error": err.Error(),
		}

		util.LogWithContext(
			p.Logger,
			slog.LevelError,
			"Failed to push notification",
			context,
			nil,
		)
	}
}

// func (s *pushService) log(context map[string]interface{}) {
// 	util.LogWithContext(
// 		s.Logger,
// 		slog.LevelInfo,
// 		"push notification sent",
// 		context,
// 		nil)
// }
