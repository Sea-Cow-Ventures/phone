package service

import (
	"sort"
	"time"

	api "github.com/twilio/twilio-go/rest/api/v2010"
	"go.uber.org/zap"

	"aidan/phone/internal/database"
	"fmt"
)

func DialNumber(agentNumber string, toNumber string) error {
	logger.Info("message")
	params := &api.CreateCallParams{}
	params.SetTo(agentNumber)
	params.SetFrom(cnf.PhoneNumber)
	params.SetUrl("https://" + cnf.UrlBasePath + "/connectAgent?toNumber=" + toNumber)

	resp, err := t.Api.CreateCall(params)
	if err != nil {
		logger.Infof("Failed to initiate call: %v", err)
		return fmt.Errorf("failed to initiate call: %w", err)
	}

	logger.Infof("Call initiated with SID: %s", *resp.Sid)
	return nil
}

func ReadAccountCallHistory() error {
	logger.Info("Reading account call history")

	calls, err := t.Api.ListCall(&api.ListCallParams{})
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
		exists, err := database.DoesCallExist(*call.Sid)
		if err != nil {
			logger.Errorf("Error reading messages: %w", zap.Error(err))
		}
		if !exists {
			err = database.InsertCall(call)
			insertedCalls++
		}
		if err != nil {
			logger.Errorf("Error inserting message: %w", zap.Error(err))
		}
	}

	logger.Info("Read account call history", zap.Int("NewCalls", insertedCalls))

	return nil
}

func OutboundAgentCall(to string) {
	params := &api.CreateCallParams{}
	params.SetTo(to)
	params.SetFrom(cnf.PhoneNumber)
	params.SetUrl(cnf.UrlBasePath + "/connectAgent")
	params.SetMachineDetection("Enable")

	resp, err := t.Api.CreateCall(params)
	logger.Info("Data", zap.Any("data", resp))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Call Status: " + *resp.Status)
		fmt.Println("Call Sid: " + *resp.Sid)
		fmt.Println("Call Direction: " + *resp.Direction)
	}
}
