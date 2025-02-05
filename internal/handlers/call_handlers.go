package handlers

import (
	"aidan/phone/internal/database"
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
