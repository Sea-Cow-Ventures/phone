package server

import (
	"aidan/phone/internal/handlers"
	"aidan/phone/internal/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.Static("/static", cnf.WebDir+"/static")
	e.Static("/static/images", cnf.WebDir+"/static/images")
	e.Static("/static/css", cnf.WebDir+"/static/css")
	e.Static("/static/js", cnf.WebDir+"/static/js")

	e.GET("/", handlers.LoginPage)
	e.GET("/login", handlers.LoginPage)
	e.GET("/home", handlers.MainPage, middleware.EnsureLoggedIn)
	e.GET("/calls", handlers.MainPage, middleware.EnsureLoggedIn)
	e.POST("/readCalls", handlers.ReadCalls, middleware.EnsureLoggedIn)
	e.POST("/markCallHandled", handlers.MarkCallHandled, middleware.EnsureLoggedIn)
	e.GET("/sms", handlers.SmsPage, middleware.EnsureLoggedIn)
	e.GET("/readMesages", handlers.ReadMessages, middleware.EnsureLoggedIn)
	e.POST("/sendMessage", handlers.SendMessage, middleware.EnsureLoggedIn)
	e.GET("/readAgents", handlers.ReadAgents, middleware.EnsureLoggedIn)
	e.POST("/addAgent", handlers.AddAgent, middleware.EnsureLoggedIn, middleware.EnsureAdmin)
	e.POST("/removeAgent", handlers.RemoveAgent, middleware.EnsureLoggedIn, middleware.EnsureAdmin)
	e.POST("/editAgent", handlers.EditAgent, middleware.EnsureLoggedIn, middleware.EnsureAdmin)
	//e.POST("/readMessagesByNumber", readMessagesByPhoneNumberHandler, middleware.EnsureLoggedIn)
	e.POST("/authenticate", handlers.Login)
	e.GET("/logout", handlers.Logout, middleware.EnsureLoggedIn)
	e.GET("/settings", handlers.SettingsPage, middleware.EnsureLoggedIn, middleware.EnsureAdmin)
	//e.POST("/sms", handlers.SendMessage, middleware.EnsureLoggedIn)
	//e.POST("/voice", voiceHandler)
	//e.POST("/welcome", welcomeHandler)
	//e.POST("/hold", holdHandler)
	e.POST("/connectAgent", handlers.ConnectAgent)
	e.POST("/dial", handlers.DialPhone, middleware.EnsureLoggedIn)
	e.POST("/voice", handlers.Voice)
	e.POST("/voice/record", handlers.VoiceRecord)
	e.POST("/fail", handlers.Fail)
	e.POST("/voice/finishRecording", handlers.FinishRecording)
	e.GET("/webpush/vapidkey", handlers.GetVAPIDKey)
	e.POST("/webpush/subscribe", handlers.SubscribeToWebpush)
	e.GET("/webpush/test", handlers.TestWebpushNotification)
}
