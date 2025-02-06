package server

import (
	"aidan/phone/internal/handlers"
	"aidan/phone/internal/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.Static("/static", "web/static")
	e.GET("/", handlers.LoginPage)
	e.GET("/login", handlers.LoginPage)
	e.GET("/home", handlers.MainPage, middleware.EnsureLoggedIn)
	e.GET("/calls", handlers.MainPage, middleware.EnsureLoggedIn)
	e.POST("/readCalls", handlers.ReadCalls, middleware.EnsureLoggedIn)
	e.GET("/sms", handlers.SmsPage, middleware.EnsureLoggedIn)
	e.GET("/readMesages", handlers.ReadMessages, middleware.EnsureLoggedIn)
	e.POST("/sendMessage", handlers.SendMessage, middleware.EnsureLoggedIn)
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
	//e.GET("/holdMusic", holdMusicHandler)
	e.POST("/dial", handlers.DialPhone, middleware.EnsureLoggedIn)
}
