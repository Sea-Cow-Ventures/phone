package database

import (
	"aidan/phone/internal/models"
)

func IsAdmin(username string) (bool, error) {
	var isAdmin bool
	err := db.Get(&isAdmin, "SELECT isAdmin FROM agents WHERE name = ?", username)
	if err != nil {
		return false, err
	}
	return isAdmin, nil
}

func GetAgentByName(name string) (*models.Agent, error) {
	agent := &models.Agent{}
	err := db.Get(agent, "SELECT id, name, number, priority, email, hashedPassword, isAdmin FROM agents WHERE name = ?", name)
	if err != nil {
		return nil, err
	}

	return agent, nil
}
