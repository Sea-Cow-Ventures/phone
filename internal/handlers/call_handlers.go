package handlers

import (
	"aidan/phone/internal/database"
	"aidan/phone/internal/models"
	"aidan/phone/internal/service"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/twilio/twilio-go/twiml"
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
		"Name":    cookie.Value,
		"IsAdmin": agent.IsAdmin,
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

func MarkCallHandled(c echo.Context) error {
	var input struct {
		CallId  int `json:"callId" validate:"required,numeric"`
		AgentId int `json:"agentId" validate:"required,numeric"`
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

	err := service.MarkCallHandled(input.CallId, input.AgentId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Failed to mark call as handled: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

func Voice(c echo.Context) error {
	call := new(models.VoiceWebhook)
	// Bind the request data to the struct
	if err := c.Bind(call); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request data")
	}

	var response []twiml.Element

	response = []twiml.Element{
		twiml.VoiceSay{
			Message:  "Welcome to Kayaking St. Augustine, we do not take bookings over the phone, leave a message and we will get back to you as soon as possible or visit www.kayaking S T augustine.com to book a tour.",
			Language: "en-US",
			Voice:    "Polly.Salli",
		},
		twiml.VoicePause{Length: "1s"},
		twiml.VoiceSay{
			Message:  "Press 1 to leave a message.",
			Language: "en-US",
			Voice:    "Polly.Salli",
		},
		twiml.VoiceGather{
			NumDigits: "1",
			Action:    "/voice/record",
			Method:    "POST",
		},
	}

	twimlResult, err := twiml.Voice(response)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.Response().Header().Set("Content-Type", "text/xml")
		c.String(http.StatusOK, twimlResult)
	}

	return nil
}

func VoiceRecord(c echo.Context) error {
	call := new(models.VoiceWebhook)
	if err := c.Bind(call); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request data")
	}

	// Check if the user pressed 1
	if call.Digits == "1" {
		response := []twiml.Element{
			twiml.VoiceSay{
				Message:  "Please leave your message after the beep. Press pound when finished.",
				Language: "en-US",
				Voice:    "Polly.Salli",
			},
			twiml.VoiceRecord{
				Action:      "/voice/finishRecording",
				Method:      "POST",
				MaxLength:   "300", // 5 minutes max
				Timeout:     "10",  // 5 seconds of silence before ending
				FinishOnKey: "#",   // Press # to end recording
				PlayBeep:    "true",
				//RecordingStatusCallback: "/voice/recording-status",
			},
		}

		twimlResult, err := twiml.Voice(response)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		c.Response().Header().Set("Content-Type", "text/xml")
		return c.String(http.StatusOK, twimlResult)
	}

	// If they didn't press 1, repeat the main menu
	response := []twiml.Element{
		twiml.VoiceSay{
			Message:  "Invalid input. Press 1 to leave a message.",
			Language: "en-US",
			Voice:    "Polly.Salli",
		},
		twiml.VoiceGather{
			NumDigits: "1",
			Action:    "/voice/finishRecording",
			Method:    "POST",
		},
	}

	twimlResult, err := twiml.Voice(response)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Set("Content-Type", "text/xml")
	return c.String(http.StatusOK, twimlResult)
}

func FinishRecording(c echo.Context) error {
	call := new(models.VoiceWebhook)
	if err := c.Bind(call); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request data")
	}

	// Log the recording details
	log.Printf("Recording URL: %s, Duration: %s seconds", call.RecordingUrl, call.RecordingDuration)

	response := []twiml.Element{
		twiml.VoiceSay{
			Message:  "Thank you for your message. We will get back to you as soon as possible. Goodbye.",
			Language: "en-US",
			Voice:    "Polly.Salli",
		},
		twiml.VoiceHangup{},
	}

	twimlResult, err := twiml.Voice(response)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Set("Content-Type", "text/xml")
	return c.String(http.StatusOK, twimlResult)
}

func Fail(c echo.Context) error {
	// Log form data
	if err := c.Request().ParseForm(); err == nil {
		formData := c.Request().Form
		logger.Info("Form Data: %+v", zap.Any("formData", formData))
	}

	// Log JSON body if present
	var bodyData map[string]interface{}
	if err := c.Bind(&bodyData); err == nil {
		logger.Info("JSON Body: %+v", zap.Any("bodyData", bodyData))
	}

	// Log cookies
	cookies := c.Cookies()
	logger.Info("Cookies: %+v", zap.Any("cookies", cookies))

	// Return a 400 Bad Request with all the logged data
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"success": false,
		"error":   "Request failed - check server logs for details",
	})
}
