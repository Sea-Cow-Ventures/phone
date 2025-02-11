package database

import (
	"aidan/phone/internal/models"
	"fmt"
)

func InsertWebpushKeys(privateKey string, publicKey string) error {
	_, err := db.Exec(`
	INSERT INTO webpushKey (publickey, privatekey)
	VALUES (?, ?)
	`, publicKey, privateKey)
	return err
}

func GetWebpushKeys() (string, string, error) {
	var privateKey, publicKey string
	err := db.QueryRow(`
	SELECT publicKey, privateKey 
	FROM webpushKey
	`).Scan(&publicKey, &privateKey)

	if err != nil {
		return "", "", fmt.Errorf("failed to get webpush keys: %w", err)
	}

	return publicKey, privateKey, nil
}

func InsertWebpushSubscription(subscription models.SubscriptionRequest, agentID int) error {
	exists, err := DoesSubscriptionExist(subscription.Endpoint)
	if err != nil {
		return err
	}

	if exists {
		_, err = db.Exec(`
			UPDATE webpushSubscriptions 
			SET 
				authKey = ?,
				p256dhKey = ?,
				lastUsed = CURRENT_TIMESTAMP
			WHERE endpoint = ?
		`,
			subscription.Keys.Auth,
			subscription.Keys.P256dh,
			subscription.Endpoint)

		if err != nil {
			return fmt.Errorf("failed to update subscription: %w", err)
		}
		return nil
	}

	_, err = db.Exec(`
		INSERT INTO webpushSubscriptions (
			agentId,
			endpoint,
			authKey,
			p256dhKey,
			enabled
		) VALUES (?, ?, ?, ?, 1)
	`,
		agentID,
		subscription.Endpoint,
		subscription.Keys.Auth,
		subscription.Keys.P256dh)

	if err != nil {
		return fmt.Errorf("failed to insert subscription: %w", err)
	}

	return nil
}

func DoesSubscriptionExist(endpoint string) (bool, error) {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM webpushSubscriptions 
			WHERE endpoint = ? 
			AND enabled = 1
		)
	`, endpoint).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check subscription existence: %w", err)
	}

	return exists, nil
}

func GetWebpushSubscriptions(agentID int) ([]models.SubscriptionRequest, error) {
	var dbSubs []models.DBSubscription
	err := db.Select(&dbSubs, `
		SELECT 
			endpoint,
			authKey,
			p256dhKey
		FROM webpushSubscriptions 
		WHERE agentId = ?
		AND enabled = 1
	`, agentID)

	if err != nil {
		return nil, fmt.Errorf("failed to get webpush subscriptions: %w", err)
	}

	subscriptions := make([]models.SubscriptionRequest, len(dbSubs))
	for i, sub := range dbSubs {
		subscriptions[i] = models.SubscriptionRequest{
			Endpoint: sub.Endpoint,
			Keys: struct {
				Auth   string `json:"auth"`
				P256dh string `json:"p256dh"`
			}{
				Auth:   sub.AuthKey,
				P256dh: sub.P256dhKey,
			},
		}
	}

	return subscriptions, nil
}
