package database

import (
	api "github.com/twilio/twilio-go/rest/api/v2010"

	"aidan/phone/internal/models"
	"fmt"
	"time"
)

func DoesMessageExist(msgSid string) (bool, error) {
	query := "SELECT COUNT(*) FROM sms WHERE messageSid = ?"

	var count int
	err := db.Get(&count, query, msgSid)
	if err != nil {
		return false, fmt.Errorf("failed to check messageSid existence: %w", err)
	}

	return count > 0, nil
}

func InsertMessageLog(msg api.ApiV2010Message) error {
	query := `
		INSERT INTO sms (
			fromNumber, toNumber, body, direction, updatedDate, price, uri,
			accountSid, mediaNumber, status, messageSid, sentDate, createdDate,
			priceUnit, apiVersion, segmentNumber, errorMessage, errorCode
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	updatedDate, err := time.Parse(time.RFC1123Z, *msg.DateUpdated)
	if err != nil {
		return fmt.Errorf("failed to parse updated date: %w", err)
	}

	sentDate, err := time.Parse(time.RFC1123Z, *msg.DateSent)
	if err != nil {
		return fmt.Errorf("failed to parse sent date: %w", err)
	}

	createdDate, err := time.Parse(time.RFC1123Z, *msg.DateCreated)
	if err != nil {
		return fmt.Errorf("failed to parse created date: %w", err)
	}

	_, err = db.Exec(
		query,
		msg.From,
		msg.To,
		msg.Body,
		msg.Direction,
		updatedDate.Format("2006-01-02 15:04:05"),
		msg.Price,
		msg.Uri,
		msg.AccountSid,
		msg.NumMedia,
		msg.Status,
		msg.Sid,
		sentDate.Format("2006-01-02 15:04:05"),
		createdDate.Format("2006-01-02 15:04:05"),
		msg.PriceUnit,
		msg.ApiVersion,
		msg.NumSegments,
		msg.ErrorMessage,
		msg.ErrorCode,
	)

	if err != nil {
		return fmt.Errorf("failed to insert message log: %w", err)
	}

	return nil
}

func GetMessagedPhoneNumbers() ([]string, error) {
	query := `
		SELECT phoneNumber
		FROM (
			SELECT fromNumber AS phoneNumber, updatedDate
			FROM sms
			WHERE fromNumber NOT IN ('+19048752208', '+19043158442')

			UNION

			SELECT toNumber AS phoneNumber, updatedDate
			FROM sms
			WHERE toNumber NOT IN ('+19048752208', '+19043158442')
		) AS combined
		GROUP BY phoneNumber
		ORDER BY MAX(updatedDate) DESC
		LIMIT 10 OFFSET 0;
	`

	var phoneNumbers []string
	err := db.Select(&phoneNumbers, query)
	if err != nil {
		return nil, err
	}

	return phoneNumbers, nil
}

func GetMessagesByPhoneNumber(phoneNumber string) ([]models.Message, error) {
	query := `
		SELECT 
			fromNumber, 
			toNumber, 
			body, 
			sentDate 
		FROM 
			sms 
		WHERE 
			fromNumber = ? 
			OR toNumber = ? 
		ORDER BY 
			sentDate ASC
	`

	var messages []models.Message
	err := db.Select(&messages, query, phoneNumber, phoneNumber)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
