package main

import (
	"fmt"
	"sync"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"go.uber.org/zap"
)

type Agent struct {
	ID             int    `db:"id" json:"id"`
	Username       string `db:"name" json:"name"`
	Password       string
	HashedPassword string `db:"hashedPassword" json:"hashedPassword"`
	Number         string `db:"number" json:"number"`
	Email          string `db:"email" json:"email"`
	Priority       int    `db:"priority" json:"priority"`
	IsAdmin        bool   `db:"isAdmin" json:"isAdmin"`
	Busy           sync.Mutex
}

type Login struct {
	Username string `db:"name" json:"name" form:"Username" validate:"required"`
	Password string `form:"Password" validate:"required"`
}

func readAgents() ([]Agent, error) {
	agents := []Agent{}
	err := db.Select(&agents, "SELECT id, name, number, priority, email, hashedPassword, isAdmin FROM agents ORDER BY id DESC")
	if err != nil {
		return nil, err
	}

	return agents, nil
}

func readAgentByName(name string) (*Agent, error) {
	agent := &Agent{}
	err := db.Get(agent, "SELECT id, name, number, priority, email, hashedPassword, isAdmin FROM agents WHERE name = ?", name)
	if err != nil {
		return nil, err
	}

	return agent, nil
}

func insertAgent(a *Agent) error {
	_, err := db.Exec("INSERT INTO agents (name, HashedPassword, number, email, priority, isAdmin) VALUES (?, ?, ?, ?, ?, ?)",
		a.Username, a.HashedPassword, a.Number, a.Email, a.Priority, a.IsAdmin)
	return err
}

func createDefaultAdmin() (*Agent, error) {
	var admin = Agent{
		Username: "admin",
		Password: generateRandomString(8),
		IsAdmin:  true,
		Priority: -1,
	}

	hashedPassword, err := hashPassword(admin.Password)
	if err != nil {
		return nil, err
	}
	admin.HashedPassword = hashedPassword

	err = insertAgent(&admin)
	return &admin, err
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
