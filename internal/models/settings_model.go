package models

type Settings struct {
	Name        string
	PhoneNumber string
	Email       string
	IsAdmin     bool
	Agents      []Agent
}
