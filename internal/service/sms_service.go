package service

import (
	"aidan/phone/internal/database"
	"fmt"
	"sort"
	"time"

	api "github.com/twilio/twilio-go/rest/api/v2010"
	"go.uber.org/zap"
)

func SendMessage(toNumber string, message string) error {
	params := &api.CreateMessageParams{}
	params.SetTo(toNumber)
	params.SetFrom(cnf.PhoneNumber)
	params.SetBody(message)

	if message == "" {
		return fmt.Errorf("message body cannot be empty")
	}

	resp, err := t.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	logger.Info("Sent message", zap.Any("response", resp))

	return nil
}

func GetAccountMessageHistory() error {
	logger.Info("Reading account message history")

	messages, err := t.Api.ListMessage(&api.ListMessageParams{})
	if err != nil {
		return err
	}

	sort.Slice(messages, func(i, j int) bool {
		createdDateI, errI := time.Parse(time.RFC1123Z, *messages[i].DateCreated)
		if errI != nil {
			logger.Errorf("Error parsing created date for index %d:%w", i, errI)
			return false
		}

		createdDateJ, errJ := time.Parse(time.RFC1123Z, *messages[j].DateCreated)
		if errJ != nil {
			logger.Errorf("Error parsing created date for index %d:%w", j, errJ)
			return false
		}

		return createdDateI.Before(createdDateJ)
	})

	var insertedMessages int

	for _, msg := range messages {
		exists, err := database.DoesMessageExist(*msg.Sid)
		if err != nil {
			return fmt.Errorf("error reading messages: %w", err)
		}

		if !exists {
			err = database.InsertMessageLog(msg)
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
