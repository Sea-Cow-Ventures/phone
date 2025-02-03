package agent

import (
	"aidan/phone/internal/server"
	"aidan/phone/pkg/util"
	"fmt"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"go.uber.org/zap"
)

func ReadAgents() ([]Agent, error) {
	agents := []Agent{}
	err := server.DB.Select(&agents, "SELECT id, name, number, priority, email, hashedPassword, isAdmin FROM agents ORDER BY id DESC")
	if err != nil {
		return nil, err
	}

	return agents, nil
}

func insertAgent(a *Agent) error {
	_, err := server.DB.Exec("INSERT INTO agents (name, HashedPassword, number, email, priority, isAdmin) VALUES (?, ?, ?, ?, ?, ?)",
		a.Username, a.HashedPassword, a.Number, a.Email, a.Priority, a.IsAdmin)
	return err
}

func CreateDefaultAdmin() (*Agent, error) {
	var admin = Agent{
		Username: "admin",
		Password: util.GenerateRandomString(8),
		IsAdmin:  true,
		Priority: -1,
	}

	hashedPassword, err := util.HashPassword(admin.Password)
	if err != nil {
		return nil, err
	}
	admin.HashedPassword = hashedPassword

	err = insertAgent(&admin)
	return &admin, err
}

func CreateAgent(username, password, email, number string, isAdmin bool) {
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		server.Logger.Error("Failed to hash password", zap.Error(err))
		return
	}

	insertAgent(&Agent{
		Username:       username,
		HashedPassword: hashedPassword,
		Email:          email,
		Number:         number,
		Priority:       0,
		IsAdmin:        isAdmin,
	})
}

func RemoveAgent(id string) {
	server.DB.Exec("DELETE FROM agents WHERE id = ?", id)
}

func EditAgent(id, username, hashedPassword, email, number string, isAdmin bool) {
	server.DB.Exec("UPDATE agents SET name = ?, hashedPassword = ?, email = ?, number = ?, isAdmin = ? WHERE id = ?",
		username, hashedPassword, email, number, isAdmin, id)
}

func IsAdmin(username string) (bool, error) {
	var isAdmin bool
	err := server.DB.Get(&isAdmin, "SELECT isAdmin FROM agents WHERE name = ?", username)
	if err != nil {
		return false, err
	}
	return isAdmin, nil
}

func IsLastAdmin(userID string) (bool, error) {
	var count int
	err := server.DB.Get(&count, "SELECT COUNT(*) FROM agents WHERE isAdmin = 1 AND id != ?", userID)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func OutboundAgentCall(to string) {
	params := &twilioApi.CreateCallParams{}
	params.SetTo(to)
	params.SetFrom(server.Cnf.PhoneNumber)
	params.SetUrl(server.Cnf.UrlBasePath + "/connectAgent")
	params.SetMachineDetection("Enable")

	resp, err := server.T.Api.CreateCall(params)
	server.Logger.Info("Data", zap.Any("data", resp))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Call Status: " + *resp.Status)
		fmt.Println("Call Sid: " + *resp.Sid)
		fmt.Println("Call Direction: " + *resp.Direction)
	}
}
