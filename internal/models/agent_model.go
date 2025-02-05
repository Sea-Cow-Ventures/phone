package models

import "sync"

type Agent struct {
	ID             int    `db:"id" json:"id"`
	Name           string `db:"name" json:"name"`
	Password       string
	HashedPassword string `db:"hashedPassword" json:"hashedPassword"`
	Number         string `db:"number" json:"number"`
	Email          string `db:"email" json:"email"`
	Priority       int    `db:"priority" json:"priority"`
	IsAdmin        bool   `db:"isAdmin" json:"isAdmin"`
	Busy           sync.Mutex
}
