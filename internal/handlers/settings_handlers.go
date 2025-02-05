package handlers

import (
	"aidan/phone/internal/models"
	"aidan/phone/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func SettingsPage(c echo.Context) error {
	cookie, err := c.Cookie("name")
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Cookie",
			Success: false,
		})
	}

	settings, err := service.GetSettings(cookie.Value)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to get settings",
			Success: false,
		})
	}

	logger.Info("Settings", zap.Any("data", settings))

	return c.Render(http.StatusOK, "settings.html", settings)
}

func RemoveAgent(c echo.Context) error {
	var input struct {
		ID int `json:"id" validate:"required,numeric"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid input format",
			Success: false,
		})
	}

	if err := c.Validate(&input); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid input format",
			Success: false,
		})
	}

	err := service.DeleteAgent(input.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to delete agent",
			Success: false,
		})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
	})
}

func AddAgent(c echo.Context) error {
	type input struct {
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Number   string `json:"number" validate:"required,e164"`
		IsAdmin  string `json:"isAdmin"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid input format",
			Success: false,
		})
	}

	if err := c.Validate(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Validation failed: " + err.Error(),
		})
	}

	agent, err := readAgentByName(input.Name)
	if agent != nil || err.Error() != "sql: no rows in result set" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Name already exists",
		})
	}

	createAgent(input.Name, input.Password, input.Email, input.Number, input.IsAdmin == "true")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
