package handlers

import (
	"aidan/phone/internal/models"
	"aidan/phone/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func SmsPage(c echo.Context) error {
	phoneNumbers, err := service.GetMessagedPhoneNumbers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to read messaged phone numbers",
		})
	}

	cookie, err := c.Cookie("name")
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Bad Cookie",
		})
	}

	agent, err := service.GetAgentByName(cookie.Value)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Bad Cookie",
		})
	}

	data := map[string]interface{}{
		"Conversations":  phoneNumbers,
		"Name":           cookie.Value,
		"IsAdmin":        agent.IsAdmin,
		"MissedCalls":    1,
		"UnreadMessages": 2,
	}

	return c.Render(http.StatusOK, "sms.html", data)
}

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

func ReadMessages(c echo.Context) error {
	phoneNumbers, err := service.GetMessagedPhoneNumbers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to read messaged phone numbers",
		})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    phoneNumbers,
	})
}
