package call

import (
	"aidan/phone/internal/server"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func ReadCallsHandler(c echo.Context) error {
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

	calls, err := readCalls(input.Page, input.Limit)
	if err != nil {
		server.Logger.Error("Failed to read calls", zap.Error(err))
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
