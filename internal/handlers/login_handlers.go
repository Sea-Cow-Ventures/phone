package handlers

import (
	"aidan/phone/internal/models"
	"aidan/phone/internal/service"
	"aidan/phone/pkg/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func LoginPage(c echo.Context) error {
	_, found := util.ReadLoginCookie(c)
	if found {
		return c.Redirect(http.StatusFound, "/home")
	}

	data := map[string]interface{}{
		"MissedCalls":    0, // Add default value for header template
		"UnreadMessages": 0, // Add default value for header template
		"Title":          "Login",
		"User":           nil, // Add nil user for header template
	}

	return c.Render(http.StatusOK, "login.html", data)
}

func Login(c echo.Context) error {
	login := new(models.Login)
	if err := c.Bind(login); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid Request Data",
		})
	}

	if err := c.Validate(login); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	err := service.Authenticate(login.Name, login.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "Unauthorized",
		})
	} else {
		util.WriteLoginCookie(c, login.Name)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success":     true,
			"redirectURL": "/home",
		})
	}
}

func Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "name"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.MaxAge = -1
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/")
}
