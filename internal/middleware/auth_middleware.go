package middleware

import (
	"aidan/phone/internal/database"
	"aidan/phone/pkg/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func EnsureLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		name, found := util.ReadLoginCookie(c)
		if !found {
			return c.Redirect(http.StatusFound, "/login")
		}
		userExists, err := database.GetAgentByName(name)
		if userExists == nil || err != nil {
			return c.Redirect(http.StatusFound, "/login")
		}

		return next(c)
	}
}

func EnsureAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		name, _ := util.ReadLoginCookie(c)
		userIsAdmin, err := database.IsAdminByName(name)
		if !userIsAdmin || err != nil {
			return c.Redirect(http.StatusFound, "/home")
		}
		return next(c)
	}
}
