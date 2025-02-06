package handlers

import (
	"aidan/phone/internal/models"
	"aidan/phone/internal/service"
	"fmt"
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
	var input struct {
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

	if err := c.Validate(&input); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation failed: " + err.Error(),
			Success: false,
		})
	}

	err := service.AddAgent(input.Name, input.Password, input.Email, input.Number, input.IsAdmin == "true")
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   fmt.Errorf("failed to add agent: %w", err).Error(),
			Success: false,
		})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
	})
}

func EditAgent(c echo.Context) error {
	var input struct {
		ID       int    `json:"id" validate:"required"`
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

	if err := c.Validate(&input); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation failed: " + err.Error(),
			Success: false,
		})
	}

	err := service.EditAgent(input.ID, input.Name, input.Password, input.Email, input.Number, input.IsAdmin == "true")
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   fmt.Errorf("failed to edit agent: %w", err).Error(),
			Success: false,
		})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
	})
}
