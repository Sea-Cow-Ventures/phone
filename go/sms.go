package main

import (
	"fmt"
	"sort"
	"time"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"go.uber.org/zap"
)

type Message struct {
	From     string `db:"fromNumber"`
	To       string `db:"toNumber"`
	Body     string `db:"body"`
	SentDate string `db:"sentDate"`
}

func readAccountMessageHistory() error {
	logger.Info("Reading account message history")

	messages, err := t.Api.ListMessage(&twilioApi.ListMessageParams{})
	if err != nil {
		return err
	}

	sort.Slice(messages, func(i, j int) bool {
		createdDateI, errI := time.Parse(time.RFC1123Z, *messages[i].DateCreated)
		if errI != nil {
			logger.Sugar().Errorf("Error parsing created date for index %d:%w", i, errI)
			return false
		}

		createdDateJ, errJ := time.Parse(time.RFC1123Z, *messages[j].DateCreated)
		if errJ != nil {
			logger.Sugar().Errorf("Error parsing created date for index %d:%w", j, errJ)
			return false
		}

		return createdDateI.Before(createdDateJ)
	})

	var insertedMessages int

	for _, msg := range messages {
		exists, err := doesMessageExist(*msg.Sid)
		if err != nil {
			return fmt.Errorf("error reading messages: %w", err)
		}

		if !exists {
			err = insertMessageLog(msg)
			insertedMessages++
		}
		if err != nil {
			return fmt.Errorf("error inserting message: %w", err)
		}

		//_, err = lookupPhoneNumber(*msg.To)
		//if err != nil {
		//	return fmt.Errorf("error looking up phone number: %w", err)
		//}

	}

	logger.Info("Read account message history", zap.Int("NewMessages", insertedMessages))

	return nil
}

func doesMessageExist(msgSid string) (bool, error) {
	query := "SELECT COUNT(*) FROM sms WHERE messageSid = ?"

	var count int
	err := db.Get(&count, query, msgSid)
	if err != nil {
		return false, fmt.Errorf("failed to check messageSid existence: %w", err)
	}

	return count > 0, nil
}

func insertMessageLog(msg twilioApi.ApiV2010Message) error {
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

func readMessagedPhoneNumbers() ([]string, error) {
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

func readMessagesByPhoneNumber(phoneNumber string) ([]Message, error) {
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

	var messages []Message
	err := db.Select(&messages, query, phoneNumber, phoneNumber)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
