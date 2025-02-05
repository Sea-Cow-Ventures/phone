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
	return c.Render(http.StatusOK, "login.html", nil)
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
