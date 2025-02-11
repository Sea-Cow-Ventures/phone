package handlers

import (
	"aidan/phone/internal/models"
	"aidan/phone/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SubscriptionRequest struct {
	Endpoint  string `json:"endpoint"`
	UserAgent string `json:"userAgent"`
	Keys      struct {
		Auth   string `json:"auth"`
		P256dh string `json:"p256dh"`
	} `json:"keys"`
}

func SubscribeToWebpush(c echo.Context) error {
	var subReq models.SubscriptionRequest
	if err := c.Bind(&subReq); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid subscription " + err.Error(),
		})
	}

	cookie, err := c.Cookie("name")
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Bad Cookie",
		})
	}

	err = service.InsertWebpushSubscription(subReq, cookie.Value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to insert subscription " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
	})
}

func GetVAPIDKey(c echo.Context) error {
	publicKey, _, err := service.GetWebpushKeys()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to get VAPID keys " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    publicKey,
	})
}

func TestWebpushNotification(c echo.Context) error {
	cookie, err := c.Cookie("name")
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Bad Cookie",
		})
	}
	err = service.SendWebpushNotification(cookie.Value, "Test message")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to send webpush notification " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
	})
}
