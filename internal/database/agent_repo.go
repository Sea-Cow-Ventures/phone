package database

import (
	"aidan/phone/internal/models"
)

func IsAdminByName(name string) (bool, error) {
	var isAdmin bool
	err := db.Get(&isAdmin, "SELECT isAdmin FROM agents WHERE name = ?", name)
	if err != nil {
		return false, err
	}
	return isAdmin, nil
}

func IsLastAdminById(id int) (bool, error) {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM agents WHERE isAdmin = 1 AND id != ?", id)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func GetAdminCount() (int, error) {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM agents WHERE isAdmin = 1")
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetAgentByName(name string) (*models.Agent, error) {
	agent := &models.Agent{}
	err := db.Get(agent, "SELECT id, name, number, priority, email, hashedPassword, isAdmin FROM agents WHERE name = ?", name)
	if err != nil {
		return nil, err
	}

	return agent, nil
}

func GetAllAgents() ([]models.Agent, error) {
	agents := []models.Agent{}
	err := db.Select(&agents, "SELECT id, name, number, priority, email, hashedPassword, isAdmin FROM agents ORDER BY id DESC")
	if err != nil {
		return nil, err
	}

	return agents, nil
}

func InsertAgent(a *models.Agent) error {
	_, err := db.Exec("INSERT INTO agents (name, HashedPassword, number, email, priority, isAdmin) VALUES (?, ?, ?, ?, ?, ?)",
		a.Name, a.HashedPassword, a.Number, a.Email, a.Priority, a.IsAdmin)
	return err
}

func DeleteAgentById(id int) error {
	_, err := db.Exec("DELETE FROM agents WHERE id = ?", id)
	return err
}

func UpdateAgentById(id int, name, hashedPassword, email, number string, isAdmin bool) error {
	_, err := db.Exec("UPDATE agents SET name = ?, hashedPassword = ?, email = ?, number = ?, isAdmin = ? WHERE id = ?",
		name, hashedPassword, email, number, isAdmin, id)
	return err
}
