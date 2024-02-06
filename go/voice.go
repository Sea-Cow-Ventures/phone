package main

import (
	"fmt"
	"sort"
	"time"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"go.uber.org/zap"
)

type VoiceWebhook struct {
	AccountSid    string `json:"AccountSid" form:"AccountSid" validate:"required"`
	ApiVersion    string `json:"ApiVersion" form:"ApiVersion" validate:"required"`
	CallSid       string `json:"CallSid" form:"CallSid" validate:"required"`
	CallStatus    string `json:"CallStatus" form:"CallStatus" validate:"required"`
	Called        string `json:"Called" form:"Called" validate:"required"`
	CalledCity    string `json:"CalledCity" form:"CalledCity" validate:"required"`
	CalledCountry string `json:"CalledCountry" form:"CalledCountry" validate:"required"`
	CalledState   string `json:"CalledState" form:"CalledState" validate:"required"`
	CalledZip     string `json:"CalledZip" form:"CalledZip" validate:"required"`
	Caller        string `json:"Caller" form:"Caller" validate:"required"`
	CallerCity    string `json:"CallerCity" form:"CallerCity" validate:"required"`
	CallerCountry string `json:"CallerCountry" form:"CallerCountry" validate:"required"`
	CallerName    string `json:"CallerName" form:"CallerName" validate:"required"`
	CallerState   string `json:"CallerState" form:"CallerState" validate:"required"`
	CallerZip     string `json:"CallerZip" form:"CallerZip" validate:"required"`
	Digits        string `json:"Digits" form:"Digits" validate:"required"`
	Direction     string `json:"Direction" form:"Direction" validate:"required"`
	FinishedOnKey string `json:"FinishedOnKey" form:"FinishedOnKey" validate:"required"`
	From          string `json:"From" form:"From" validate:"required"`
	FromCity      string `json:"FromCity" form:"FromCity" validate:"required"`
	FromCountry   string `json:"FromCountry" form:"FromCountry" validate:"required"`
	FromState     string `json:"FromState" form:"FromState" validate:"required"`
	FromZip       string `json:"FromZip" form:"FromZip" validate:"required"`
	To            string `json:"To" form:"To" validate:"required"`
	ToCity        string `json:"ToCity" form:"ToCity" validate:"required"`
	ToCountry     string `json:"ToCountry" form:"ToCountry" validate:"required"`
	ToState       string `json:"ToState" form:"ToState" validate:"required"`
	ToZip         string `json:"ToZip" form:"ToZip" validate:"required"`
	Msg           string `json:"msg" form:"msg" validate:"required"`
}

func readAccountCallHistory() error {
	logger.Info("Reading account call history")

	calls, err := t.Api.ListCall(&twilioApi.ListCallParams{})
	if err != nil {
		return err
	}

	sort.Slice(calls, func(i, j int) bool {
		createdDateI, errI := time.Parse(time.RFC1123Z, *calls[i].DateCreated)
		if errI != nil {
			fmt.Println("Error parsing created date for index", i, ":", errI)
			return false
		}

		createdDateJ, errJ := time.Parse(time.RFC1123Z, *calls[j].DateCreated)
		if errJ != nil {
			fmt.Println("Error parsing created date for index", j, ":", errJ)
			return false
		}

		return createdDateI.Before(createdDateJ)
	})

	var insertedCalls int

	for _, call := range calls {
		exists, err := doesCallExist(*call.Sid)
		if err != nil {
			fmt.Println("Error reading messages:", err)
		}
		if !exists {
			err = insertCallLog(call)
			insertedCalls++
		}
		if err != nil {
			fmt.Println("Error inserting message:", err)
		}
	}

	logger.Info("Read account call history", zap.Int("NewCalls", insertedCalls))

	return nil
}

func doesCallExist(callSid string) (bool, error) {
	query := "SELECT COUNT(*) FROM calls WHERE callSid = ?"

	var count int
	err := db.Get(&count, query, callSid)
	if err != nil {
		return false, fmt.Errorf("failed to check callSid existence: %w", err)
	}

	return count > 0, nil
}

func insertCallLog(call twilioApi.ApiV2010Call) error {
	query := `
		INSERT INTO calls (
			fromNumber, toNumber, direction, updatedDate, price, uri,
			accountSid, status, callSid, sentDate, createdDate,
			priceUnit, apiVersion, parentCallSid,
			toFormatted, fromFormatted, phoneNumberSid, answeredBy,
			forwardedFrom, groupSid, callerName, queueTime, trunkSid
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	updatedDate, err := time.Parse(time.RFC1123Z, *call.DateUpdated)
	if err != nil {
		return fmt.Errorf("failed to parse updated date: %w", err)
	}

	sentDate, err := time.Parse(time.RFC1123Z, *call.StartTime)
	if err != nil {
		return fmt.Errorf("failed to parse sent date: %w", err)
	}

	createdDate, err := time.Parse(time.RFC1123Z, *call.DateCreated)
	if err != nil {
		return fmt.Errorf("failed to parse created date: %w", err)
	}

	_, err = db.Exec(
		query,
		call.From,
		call.To,
		call.Direction,
		updatedDate.Format("2006-01-02 15:04:05"),
		call.Price,
		call.Uri,
		call.AccountSid,
		call.Status,
		call.Sid,
		sentDate.Format("2006-01-02 15:04:05"), // Format time for MySQL DATETIME
		createdDate.Format("2006-01-02 15:04:05"), // Format time for MySQL DATETIME
		call.PriceUnit,
		call.ApiVersion,
		call.ParentCallSid,
		call.ToFormatted,
		call.FromFormatted,
		call.PhoneNumberSid,
		call.AnsweredBy,
		call.ForwardedFrom,
		call.GroupSid,
		call.CallerName,
		call.QueueTime,
		call.TrunkSid,
	)

	if err != nil {
		return fmt.Errorf("failed to insert call log: %w", err)
	}

	return nil
}
