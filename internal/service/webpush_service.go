package service

import (
	"aidan/phone/internal/database"
	"aidan/phone/internal/models"
	"io"
	"net/http"

	"github.com/SherClockHolmes/webpush-go"
	"go.uber.org/zap"
)

func GetWebpushKeys() (string, string, error) {
	publicKey, privateKey, err := database.GetWebpushKeys()
	if err != nil {
		return "", "", err
	}
	return publicKey, privateKey, nil
}

func GenerateWebpushKeys() (string, string, error) {
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		return "", "", err
	}

	err = database.InsertWebpushKeys(privateKey, publicKey)
	if err != nil {
		return "", "", err
	}

	return privateKey, publicKey, nil
}

func InsertWebpushSubscription(subscription models.SubscriptionRequest, name string) error {
	agent, err := database.GetAgentByName(name)
	if err != nil {
		return err
	}

	err = database.InsertWebpushSubscription(subscription, agent.ID)
	if err != nil {
		return err
	}
	return nil
}

func SendWebpushNotification(name, message string) error {
	agent, err := database.GetAgentByName(name)
	if err != nil {
		return err
	}

	subscriptions, err := database.GetWebpushSubscriptions(agent.ID)
	if err != nil {
		return err
	}

	vapidPublicKey, vapidPrivateKey, err := GetWebpushKeys()
	if err != nil {
		return err
	}

	for _, subscription := range subscriptions {
		var resp *http.Response
		resp, err = webpush.SendNotification([]byte(message), &webpush.Subscription{
			Endpoint: subscription.Endpoint,
			Keys: webpush.Keys{
				Auth:   subscription.Keys.Auth,
				P256dh: subscription.Keys.P256dh,
			},
		}, &webpush.Options{
			Subscriber:      "aidan@kayakingstaugustine.com",
			VAPIDPublicKey:  vapidPublicKey,
			VAPIDPrivateKey: vapidPrivateKey,
			TTL:             30,
		})
		if err != nil {
			return err
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		logger.Info("Sent webpush notification", zap.Int("response code", resp.StatusCode), zap.Any("response", body))
	}

	return nil
}
