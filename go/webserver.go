package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"text/template"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	tValidatorClient "github.com/twilio/twilio-go/client"
	"github.com/twilio/twilio-go/twiml"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/bcrypt"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

var upgrader = websocket.Upgrader{}

func liveReloadHandler(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// Watch the directory containing your views
	err = watcher.Add("go/views")
	if err != nil {
		return err
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				err = conn.WriteMessage(websocket.TextMessage, []byte("reload"))
				if err != nil {
					return err
				}
			}
		case err := <-watcher.Errors:
			return err
		}
	}
}

func initWebserver() {
	e.HideBanner = true
	e.HidePort = true

	e.Use(recoverMiddleware())
	e.Use(loggerMiddleware())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))
	e.Validator = &CustomValidator{validator: validator.New()}

	// Register the template renderer with the toJsonm function
	e.Renderer = &Template{
		templates: template.Must(template.New("").Funcs(template.FuncMap{
			"toJSON": toJSON,
		}).ParseGlob("views/*.html")),
	}

	// Define routes
	e.Static("/", "views")
	e.GET("/", loginHandler)
	e.GET("/home", homeHandler, isLoggedIn)
	e.GET("/smsLog", smsLogHandler, isLoggedIn)
	e.GET("/readMessagedPhoneNumbers", readMessagedPhoneNumbersHandler, isLoggedIn)
	e.POST("/sendMessage", sendMessageHandler, isLoggedIn)
	e.POST("/readMessageHistory", readMessagesByPhoneNumberHandler, isLoggedIn)
	e.POST("/signin", signinHandler)
	e.GET("/logout", logoutHandler, isLoggedIn)
	e.GET("/settings", settingsHandler, isLoggedIn)
	e.POST("/sms", smsHandler)
	e.POST("/voice", voiceHandler)
	e.POST("/welcome", welcomeHandler)
	e.POST("/hold", holdHandler)
	e.POST("/connectAgent", connectAgentHandler)
	e.GET("/holdMusic", holdMusicHandler)
	e.GET("/livereload", liveReloadHandler)
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

func loggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			fields := []zapcore.Field{
				zap.String("request", fmt.Sprintf("%s %s", req.Method, req.RequestURI)),
				zap.Int("status", res.Status),
				zap.String("host", req.Host),
				zap.String("remote_ip", c.RealIP()),
				zap.String("user_agent", req.UserAgent()),
				zap.String("latency", time.Since(start).String()),
				zap.Int64("size", res.Size),
			}

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}
			fields = append(fields, zap.String("request_id", id))

			n := res.Status
			switch {
			case n >= 500:
				body, err := io.ReadAll(c.Request().Body)
				if err != nil {
					logger.Error("Failed to read request body", zap.Error(err))
					return err
				}
				c.Request().Body = io.NopCloser(bytes.NewReader(body))

				fields = append(fields, zap.String("request_body", string(body)))
				fields = append(fields, zap.Error(err))

				logger.With(zap.Error(err)).Error("Webserver error", fields...)
				//emailErr := email.SendErrorEmail(
				//	config.MailServer,
				//	config.ServiceName,
				//	err,
				//	config.EmailRecipients,
				//	config.EmailCC,
				//	config.EmailBCC,
				//	config.FromEmail,
				//)
				//if emailErr != nil {
				//	logger.Error("Sending error email", zap.Error(emailErr))
				//}
			case n >= 400:
				body, err := io.ReadAll(c.Request().Body)
				if err != nil {
					logger.Error("Failed to read request body", zap.Error(err))
					return err
				}

				fields = append(fields, zap.String("request_body", string(body)))

				logger.With(zap.Error(err)).Warn("Webserver client error", fields...)
			case n >= 300:
				logger.Info("Webserver redirection", fields...)
			default:
				logger.Info("Webserver success", fields...)
			}

			return nil
		}
	}
}

func recoverMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					logger.Error("Webserver error", zap.Error(err), zap.Stack("stack"))
				} else {
					logger.Error("Webserver error", zap.Error(err), zap.Stack("stack"))
					//emailErr := email.SendErrorEmail(
					//	config.MailServer,
					//	config.ServiceName,
					//	err,
					//	config.EmailRecipients,
					//	config.EmailCC,
					//	config.EmailBCC,
					//	config.FromEmail,
					//)
					//if emailErr != nil {
					//	logger.Error("Sending error email", zap.Error(emailErr))
					//}

					return c.JSON(http.StatusInternalServerError, ErrorResponse{err.Error(), false})
				}
			}
			return next(c)
		}
	}
}

func isLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, found := readLoginCookie(c)
		if !found {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error":   "Unauthorized",
			})
		}
		return next(c)
	}
}

func smsHandler(c echo.Context) error {

	//logger.Info("Received sms", zap.Any("msg", data)

	bytes, _ := io.ReadAll(c.Request().Body)

	logger.Info("data", zap.Any("data", bytes))

	tValidator := tValidatorClient.NewRequestValidator(cnf.TwilioPass)

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

	valid := tValidator.ValidateBody("https://"+cnf.UrlBasePath+"/sms", bytes, signature)

	logger.Info("Validation", zap.String("url", "https://"+cnf.UrlBasePath+"/sms"), zap.Bool("Validation", valid))

	if !valid {
		logger.Error("Request failed twilio validation", zap.Error(err))
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

	file, err := os.Open(cnf.HoldMusicPath)
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

func welcomeHandler(c echo.Context) error {
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
		logger.Fatal("Failed to read agents", zap.Error(err))
	}

	for _, agent := range agents {
		if agent.Busy.TryLock() {
			outboundAgentCall(agent.Number)
		}
	}

	return nil
}

func connectAgentHandler(c echo.Context) error {
	fmt.Println("GET Data:")
	for key, value := range c.QueryParams() {
		fmt.Printf("%s: %s\n", key, value)
	}

	// Capture and print POST data
	fmt.Println("POST Data:")
	if c.Request().Method == http.MethodPost {
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			fmt.Println("Error reading request body:", err)
		} else {
			fmt.Println(string(body))
		}
	}

	play := []twiml.Element{
		twiml.VoiceDial{InnerElements: []twiml.Element{
			twiml.VoiceQueue{Name: "rep"},
		},
		},
	}
	twimlResult, err := twiml.Voice(play)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.Response().Header().Set("Content-Type", "text/xml")
		c.String(http.StatusOK, twimlResult)
	}

	return nil
}

func homeHandler(c echo.Context) error {
	cookie, err := c.Cookie("username")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Bad Cookie",
		})
	}

	data := map[string]interface{}{
		"Username":       cookie.Value,
		"MissedCalls":    1,
		"UnreadMessages": 2,
	}

	return c.Render(http.StatusOK, "home.html", data)
}

func logoutHandler(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "username"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.MaxAge = -1
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/")
}

func settingsHandler(c echo.Context) error {
	cookie, err := c.Cookie("username")
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
			"error":   "Bad Username",
		})
	}

	allAgentData, err := readAgents()
	if err != nil {
		logger.Error("Failed to read agents", zap.Error(err))
	}

	type outputAgent struct {
		ID       int
		Username string
		Number   string
		Email    string
		IsAdmin  bool
	}

	agents := []outputAgent{}
	for _, agent := range allAgentData {
		agents = append(agents, outputAgent{
			ID:       agent.ID,
			Username: agent.Username,
			Number:   agent.Number,
			Email:    agent.Email,
			IsAdmin:  agent.IsAdmin,
		})
	}

	data := map[string]interface{}{
		"Username":       agent.Username,
		"PhoneNumber":    agent.Number,
		"Email":          agent.Email,
		"IsAdmin":        agent.IsAdmin,
		"Agents":         agents,
		"MissedCalls":    1,
		"UnreadMessages": 2,
	}

	logger.Info("Settings", zap.Any("data", data))

	return c.Render(http.StatusOK, "settings.html", data)
}

func smsLogHandler(c echo.Context) error {
	phoneNumbers, err := readMessagedPhoneNumbers()
	if err != nil {
		logger.Error("Failed to read messaged phone numbers", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to read messaged phone numbers",
		})
	}

	cookie, err := c.Cookie("username")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Bad Cookie",
		})
	}

	data := map[string]interface{}{
		"Conversations":  phoneNumbers,
		"Username":       cookie.Value,
		"MissedCalls":    1,
		"UnreadMessages": 2,
	}

	return c.Render(http.StatusOK, "smsLog.html", data)
}

func readMessagedPhoneNumbersHandler(c echo.Context) error {
	phoneNumbers, err := readMessagedPhoneNumbers()
	if err != nil {
		logger.Error("Failed to read messaged phone numbers", zap.Error(err))
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
	logger.Info("Messages", zap.Any("messages", messages))
	if err != nil {
		logger.Error("Failed to read messages by phone number", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to read messages by phone number",
		})
	}

	return c.JSON(http.StatusOK, messages)
}

func sendMessageHandler(c echo.Context) error {
	toNumber := c.FormValue("toNumber")
	message := c.FormValue("message")

	err := sendMessage(toNumber, message)
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

func loginHandler(c echo.Context) error {
	_, found := readLoginCookie(c)
	if found {
		return c.Redirect(http.StatusFound, "/home")
	}
	return c.Render(http.StatusOK, "login.html", nil)
}

func signinHandler(c echo.Context) error {
	login := new(Login)
	if err := c.Bind(login); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid Request Data",
		})
	}

	if err := c.Validate(login); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	agent, err := readAgentByName(login.Username)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(agent.HashedPassword), []byte(login.Password))

	if err == nil {
		writeLoginCookie(c, agent.Username)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success":     true,
			"redirectURL": "/home",
		})
	}

	return c.JSON(http.StatusUnauthorized, map[string]interface{}{
		"success": false,
		"error":   "Unauthorized",
	})
}
