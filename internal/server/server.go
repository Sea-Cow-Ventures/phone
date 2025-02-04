package server

import (
	"aidan/phone/pkg/util"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"text/template"

	"aidan/phone/internal/config"
	"aidan/phone/internal/database"
	"aidan/phone/internal/log"
	customMiddleware "aidan/phone/internal/middleware"
	"aidan/phone/internal/models"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	tValidatorClient "github.com/twilio/twilio-go/client"

	"github.com/twilio/twilio-go"
	"github.com/twilio/twilio-go/twiml"
	"go.uber.org/zap"
)

var (
	Cnf    config.AppConfig
	Logger *zap.SugaredLogger
	DB     *sqlx.DB
	T      *twilio.RestClient
)

func init() {
	Cnf = config.GetConfig()
	Logger = log.GetLogger()
	DB = database.GetDb()
}

func Start() {
	e := echo.New()
	initServer(e)

	// Initial template loading
	loadTemplates()

	go startLiveReloadWatcher()

	e.HideBanner = true
	e.HidePort = true

	e.Use(customMiddleware.Recover())
	e.Use(customMiddleware.Log())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))
	e.Validator = &models.Validator{Validator: validator.New()}

	// Define routes

}

func smsHandler(c echo.Context) error {

	//Logger.Info("Received sms", zap.Any("msg", data)

	bytes, _ := io.ReadAll(c.Request().Body)

	Logger.Info("data", zap.Any("data", bytes))

	tValidator := tValidatorClient.NewRequestValidator(Cnf.TwilioPass)

	signature := c.Request().Header.Get("X-Twilio-Signature")

	parsedData, err := url.ParseQuery(string(bytes))

	resultMap := make(map[string]string)
	for key, values := range parsedData {
		if len(values) > 0 {
			resultMap[key] = values[0]
		}
	}

	jsonString, err := json.MarshalIndent(parsedData, "", "    ")
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Print the pretty-printed JSON string
	fmt.Println(string(jsonString))

	//bodyBytes, err := io.ReadAll(c.Request().Body)
	//if err != nil {
	//	return err
	//}

	valid := tValidator.ValidateBody("https://"+Cnf.UrlBasePath+"/sms", bytes, signature)

	Logger.Info("Validation", zap.String("url", "https://"+Cnf.UrlBasePath+"/sms"), zap.Bool("Validation", valid))

	if !valid {
		Logger.Error("Request failed twilio validation", zap.Error(err))
	}

	//blocks, err := getBlockedNumbers()
	//if err != nil {
	//	logger.Error("Failed to get blocked phone numbers", zap.Error(err))
	//}

	var message = new(twiml.MessagingMessage)
	message.Body = "Bingo Bango"

	//for _, block := range blocks {
	//	if data["From"] == block.PhoneNumber {
	//		message.Body = "STOP"
	//	}
	//}

	twimlResult, err := twiml.Messages([]twiml.Element{message})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		//c.Header("Content-Type", "text/xml")
		c.Response().Header().Set("Content-Type", "text/xml")
		c.String(http.StatusOK, twimlResult)
	}
	return nil
}

func voiceHandler(c echo.Context) error {
	//play := &twiml.VoicePlay{
	//	Url: tunnel + "/testAudio",
	//}

	voiceBody := []twiml.Element{
		twiml.VoicePause{Length: "3"},
		twiml.VoiceGather{
			Action:    "/welcome",
			Method:    "POST",
			NumDigits: "1",
			Input:     "dtmf",
			Timeout:   "10",
			InnerElements: []twiml.Element{twiml.VoiceSay{
				Message:  "Welcome to Kayaking St. Augustine! Press 1 to speak to a representative.",
				Language: "en-US",
				Voice:    "Polly.Salli",
			},
			},
		},
		twiml.VoicePause{Length: "5"},
	}

	/*stream := twiml.VoiceStream{
		Name: "Voice Stream Handler",
		Url:  "wss://fully-lenient-grouse.ngrok-free.app/voiceStream",
	}
	*/

	twimlResult, err := twiml.Voice(voiceBody)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.Response().Header().Set("Content-Type", "text/xml")
		c.String(http.StatusOK, twimlResult)
	}

	fmt.Printf("%+v\n", twimlResult)

	return nil
	/*
	   c.Response().Header().Set("Content-Type", "text/xml")

	   c.String(http.StatusOK, fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><Response><Connect><Stream name="Voice Stream" url="wss://fully-lenient-grouse.ngrok-free.app/voiceStream" /></Connect></Response>`))

	   return nil
	*/
}

func holdMusicHandler(c echo.Context) error {
	fmt.Println("audio")

	file, err := os.Open(Cnf.HoldMusicPath)
	if err != nil {
		return err
	}
	defer file.Close()

	c.Response().Header().Set(echo.HeaderContentType, "audio/ogg")

	bufferSize := 1024
	buffer := make([]byte, bufferSize)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		_, err = c.Response().Write(buffer[:n])
		if err != nil {
			return err
		}
	}

	return nil
}

/*func welcomeHandler(c echo.Context) error {
	call := new(VoiceWebhook)
	// Bind the request data to the struct
	if err := c.Bind(call); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request data")
	}

	// Validate the request data
	//if err := c.Validate(call); err != nil {
	//	return c.String(http.StatusBadRequest, fmt.Sprintf("Validation error: %s", err.Error()))
	//}

	fmt.Println(call.Digits)

	var response []twiml.Element

	if call.Digits == "1" {
		response = []twiml.Element{
			twiml.VoiceSay{
				Message:  "Connecting you with someone who can help.",
				Language: "en-US",
				Voice:    "Polly.Salli",
			},
			twiml.VoiceEnqueue{
				Name:          "Rep",
				WaitUrlMethod: "POST",
				WaitUrl:       "/hold",
			},
			//twiml.VoiceDial{
			//	Number: "+18157017775",
			//},
		}
	} else {
		response = []twiml.Element{
			twiml.VoiceGather{
				Action:    "/welcome",
				Method:    "POST",
				NumDigits: "1",
				Input:     "dtmf",
				Timeout:   "10",
				InnerElements: []twiml.Element{twiml.VoiceSay{
					Message:  "Sorry that response was invalid. Please try again.",
					Language: "en-US",
					Voice:    "Polly.Salli",
				},
				},
			},
		}
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

func holdHandler(c echo.Context) error {
	play := []twiml.Element{twiml.VoicePlay{
		Url: "/holdMusic",
	}}
	twimlResult, err := twiml.Voice(play)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.Response().Header().Set("Content-Type", "text/xml")
		c.String(http.StatusOK, twimlResult)
	}

	agents, err := readAgents()
	if err != nil {
		Logger.Fatal("Failed to read agents", zap.Error(err))
	}

	for _, agent := range agents {
		if agent.Busy.TryLock() {
			outboundAgentCall(agent.Number)
		}
	}

	return nil
}*/

func connectAgentHandler(c echo.Context) error {
	toNumber := c.QueryParam("toNumber")
	if toNumber == "" {
		return c.String(http.StatusBadRequest, "Missing 'toNumber' parameter")
	}

	voiceBody := []twiml.Element{
		twiml.VoiceDial{Number: toNumber},
	}

	twimlResult, err := twiml.Voice(voiceBody)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Set("Content-Type", "text/xml")
	return c.String(http.StatusOK, twimlResult)
}

func dialHandler(c echo.Context) error {
	type dialInput struct {
		PhoneNumber string `json:"phoneNumber" validate:"required"`
	}
	input := new(dialInput)

	if err := c.Bind(&input); err != nil {
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

	cookie, err := c.Cookie("name")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Bad Cookie",
		})
	}

	agent, err := readAgentByName(cookie.Value)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Bad Name",
		})
	}

	dialNumber(agent.Number, input.PhoneNumber)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

func logoutHandler(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "name"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.MaxAge = -1
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/")
}

//func notificationsHandler(c echo.Context) error {
//}

func smsLogHandler(c echo.Context) error {
	phoneNumbers, err := readMessagedPhoneNumbers()
	if err != nil {
		Logger.Error("Failed to read messaged phone numbers", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to read messaged phone numbers",
		})
	}

	cookie, err := c.Cookie("name")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Bad Cookie",
		})
	}

	data := map[string]interface{}{
		"Conversations":  phoneNumbers,
		"Name":           cookie.Value,
		"MissedCalls":    1,
		"UnreadMessages": 2,
	}

	return c.Render(http.StatusOK, "smsLog.html", data)
}

func readMessagedPhoneNumbersHandler(c echo.Context) error {
	phoneNumbers, err := readMessagedPhoneNumbers()
	if err != nil {
		Logger.Error("Failed to read messaged phone numbers", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to read messaged phone numbers",
		})
	}

	return c.JSON(http.StatusOK, phoneNumbers)
}

func readMessagesByPhoneNumberHandler(c echo.Context) error {
	phoneNumber := c.FormValue("phoneNumber")
	messages, err := readMessagesByPhoneNumber(phoneNumber)
	Logger.Info("Messages", zap.Any("messages", messages))
	if err != nil {
		Logger.Error("Failed to read messages by phone number", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to read messages by phone number",
		})
	}

	return c.JSON(http.StatusOK, messages)
}

func addUserHandler(c echo.Context) error {
	type addUserInput struct {
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Number   string `json:"number" validate:"required,e164"`
		IsAdmin  string `json:"isAdmin"`
	}

	input := new(addUserInput)
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

func editUserHandler(c echo.Context) error {
	type addUserInput struct {
		UserID   string `json:"userId" validate:"required"`
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Number   string `json:"number" validate:"required,e164"`
		IsAdmin  string `json:"isAdmin"`
	}

	input := new(addUserInput)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid input format",
		})
	}

	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to hash password",
		})
	}

	editAgent(input.UserID, input.name, hashedPassword, input.Email, input.Number, input.IsAdmin == "true")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

func loadTemplates() {
	parsedTemplates, err := template.New("").Funcs(template.FuncMap{
		"toJSON": util.ToJSON,
	}).ParseGlob("web/templates/*.html")
	if err != nil {
		Logger.Fatal("Failed to parse templates", zap.Error(err))
	}

	e.Renderer = &Template{
		templates: parsedTemplates,
	}
}

func startLiveReloadWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		Logger.Fatal("Failed to create fsnotify watcher", zap.Error(err))
	}
	defer watcher.Close()

	// Watch the directory containing your views
	err = watcher.Add("web/templates")
	if err != nil {
		Logger.Fatal("Failed to add fsnotify watcher", zap.Error(err))
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				Logger.Info("Live reload triggered")
				loadTemplates()
			}
		case err := <-watcher.Errors:
			Logger.Error("Error in fsnotify watcher", zap.Error(err))
		}
	}
}
