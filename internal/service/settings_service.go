package service

import (
	"aidan/phone/internal/database"
	"aidan/phone/internal/models"
	"aidan/phone/pkg/util"
	"errors"

	"go.uber.org/zap"
)

func CreateAgent(name, password, email, number string, isAdmin bool) {
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		logger.Error("Failed to hash password", zap.Error(err))
		return
	}

	database.InsertAgent(&models.Agent{
		Name:           name,
		HashedPassword: hashedPassword,
		Email:          email,
		Number:         number,
		Priority:       0,
		IsAdmin:        isAdmin,
	})
}

func CreateDefaultAdmin() (*models.Agent, error) {
	var admin = models.Agent{
		Name:     "admin",
		Password: util.GenerateRandomString(8),
		IsAdmin:  true,
		Priority: -1,
	}

	hashedPassword, err := util.HashPassword(admin.Password)
	if err != nil {
		return nil, err
	}
	admin.HashedPassword = hashedPassword

	err = database.InsertAgent(&admin)
	return &admin, err
}

func GetAgentByName(name string) (*models.Agent, error) {
	return database.GetAgentByName(name)
}

func GetSettings(name string) (*models.Settings, error) {
	agent, err := database.GetAgentByName(name)
	if err != nil {
		return nil, err
	}

	allAgents, err := database.GetAllAgents()
	if err != nil {
		return nil, err
	}

	settings := models.Settings{
		Name:        agent.Name,
		PhoneNumber: agent.Number,
		Email:       agent.Email,
		IsAdmin:     agent.IsAdmin,
		Agents:      allAgents,
	}

	return &settings, nil
}

func AddAgent(name, password, email, number string, isAdmin bool) error {
	agent, err := database.GetAgentByName(name)
	if agent != nil || err.Error() != "sql: no rows in result set" {
		return errors.New("agent already exists")
	}

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return err
	}

	err = database.InsertAgent(&models.Agent{
		Name:           name,
		HashedPassword: hashedPassword,
		Email:          email,
		Number:         number,
		Priority:       0,
		IsAdmin:        isAdmin,
	})

	return err
}

func DeleteAgent(id int) error {
	isLast, err := database.IsLastAdminById(id)
	if err != nil {
		return err
	}

	if isLast {
		return errors.New("Cannot delete last admin user")
	}

	err = database.DeleteAgentById(id)
	if err != nil {
		return err
	}

	return nil
}

func EditAgent(id int, name, password, email, number string, isAdmin bool) error {
	agent, err := database.GetAgentByName(name)
	if err != nil || agent.ID != id {
		return errors.New("agent not found")
	}

	if agent.Name == name &&
		agent.Email == email &&
		agent.Number == number &&
		agent.IsAdmin == isAdmin {
		return nil
	}

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return err
	}

	return database.UpdateAgentById(id, name, hashedPassword, email, number, isAdmin)
}
