package service

import (
	"aidan/phone/internal/database"

	"golang.org/x/crypto/bcrypt"
)

func Authenticate(name string, password string) error {
	agent, err := database.GetAgentByName(name)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(agent.HashedPassword), []byte(password))

	if err != nil {
		return err
	}

	return nil
}
