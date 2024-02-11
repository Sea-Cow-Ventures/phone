package main

type User struct {
	Username     string `form:"username" validate:"required"`
	Password     string `form:"Password" validate:"required"`
	PasswordHash string `db:"passwordHash"`
}
