package handlers

import (
	"aidan/phone/internal/database"
	"aidan/phone/internal/models"
	"aidan/phone/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func ReadCalls(c echo.Context) error {
	type readCallsInput struct {
		Page  int `json:"page" validate:"required"`
		Limit int `json:"limit" validate:"required"`
	}

	input := new(readCallsInput)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid input format",
		})
	}

	if err := c.Validate(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Validation failed: " + err.Error(),
		})
	}

	calls, err := database.ReadCalls(input.Page, input.Limit)
	if err != nil {
		logger.Error("Failed to read calls", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to read calls",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    calls,
	})
}

func MainPage(c echo.Context) error {
	cookie, err := c.Cookie("name")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Bad Cookie",
		})
	}

	agent, err := service.GetAgentByName(cookie.Value)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Bad Cookie",
		})
	}

	data := map[string]interface{}{
		"Name":           cookie.Value,
		"IsAdmin":        agent.IsAdmin,
		"MissedCalls":    1,
		"UnreadMessages": 2,
	}

	return c.Render(http.StatusOK, "home.html", data)
}

func DialPhone(c echo.Context) error {
	var input struct {
		PhoneNumber string `json:"phoneNumber" validate:"required"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid input format",
		})
	}

	if err := c.Validate(input); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Validation failed: " + err.Error(),
		})
	}

	cookie, err := c.Cookie("name")
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Bad Cookie",
		})
	}

	err = service.DialPhone(input.PhoneNumber, cookie.Value)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Failed to dial phone: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

func ConnectAgent(c echo.Context) error {
	toNumber := c.QueryParam("toNumber")
	if toNumber == "" {
		return c.String(http.StatusBadRequest, "Missing 'toNumber' parameter")
	}

	twiml, err := service.ConnectAgent(toNumber)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, twiml)
}
