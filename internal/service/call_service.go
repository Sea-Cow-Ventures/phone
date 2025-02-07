package service

import (
	"sort"
	"time"

	api "github.com/twilio/twilio-go/rest/api/v2010"
	"go.uber.org/zap"

	"aidan/phone/internal/database"
	"fmt"

	"github.com/twilio/twilio-go/twiml"
)

func GetAccountCallHistory() error {
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
			logger.Sugar().Errorf("Error reading messages: %w", zap.Error(err))
		}
		if !exists {
			err = database.InsertCall(call)
			insertedCalls++
		}
		if err != nil {
			logger.Sugar().Errorf("Error inserting message: %w", zap.Error(err))
		}
	}

	logger.Info("Read account call history", zap.Int("NewCalls", insertedCalls))

	return nil
}

func DialPhone(toPhoneNumber, agentName string) error {
	agent, err := database.GetAgentByName(agentName)
	if err != nil {
		return fmt.Errorf("error reading agent: %w", err)
	}

	params := &api.CreateCallParams{}
	params.SetTo(agent.Number)
	params.SetFrom(cnf.PhoneNumber)
	params.SetUrl("https://" + cnf.UrlBasePath + "/connectAgent?toNumber=" + toPhoneNumber)

	resp, err := t.Api.CreateCall(params)
	if err != nil {
		return fmt.Errorf("failed to initiate call: %w", err)
	}

	logger.Sugar().Infof("Call initiated with SID: %s", *resp.Sid)
	return nil
}

func ConnectAgent(toNumber string) (string, error) {
	voiceBody := []twiml.Element{
		twiml.VoiceDial{Number: toNumber},
	}

	twimlResult, err := twiml.Voice(voiceBody)
	if err != nil {
		return "", fmt.Errorf("failed to generate twiml: %w", err)
	}

	return twimlResult, nil
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

func MarkCallHandled(callSid string) error {
	err := database.MarkCallHandled(callSid)
	if err != nil {
		return fmt.Errorf("error marking call as handled: %w", err)
	}
	return nil
}
