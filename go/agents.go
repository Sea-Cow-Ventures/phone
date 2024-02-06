package main

import (
	"fmt"
	"sync"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"go.uber.org/zap"
)

type Agent struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Number   string `db:"number" json:"number"`
	Priority int    `db:"priority" json:"priority"`
	Busy     sync.Mutex
}

func readAgents() ([]Agent, error) {
	agents := []Agent{}
	err := db.Select(&agents, "SELECT id, name, number, priority FROM agents ORDER BY priority ASC")
	if err != nil {
		//logger.Fatal(err)
		return nil, err
	}

	return agents, nil
}

func outboundAgentCall(to string) {
	params := &twilioApi.CreateCallParams{}
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
