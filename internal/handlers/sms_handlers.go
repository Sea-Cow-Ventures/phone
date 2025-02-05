package handlers

import (
	"aidan/phone/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func SendMessage(c echo.Context) error {
	toNumber := c.FormValue("toNumber")
	message := c.FormValue("message")

	err := service.SendMessage(toNumber, message)
	if err != nil {
		logger.Error("Failed to send message to "+toNumber, zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to send message to " + toNumber,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
