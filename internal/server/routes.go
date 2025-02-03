package server

import (
	"aidan/phone/internal/handlers"
	"aidan/phone/internal/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.Static("/static", "web/static")
	e.GET("/", loginHandler)
	e.GET("/login", loginHandler)
	e.GET("/home", homeHandler, middleware.EnsureLoggedIn)
	e.GET("/calls", homeHandler, middleware.EnsureLoggedIn)
	e.POST("/readCalls", handlers.ReadCalls, middleware.EnsureLoggedIn)
	e.GET("/smsLog", smsLogHandler, middleware.EnsureLoggedIn)
	e.GET("/readMessagedPhoneNumbers", readMessagedPhoneNumbersHandler, middleware.EnsureLoggedIn)
	e.POST("/sendMessage", sendMessageHandler, middleware.EnsureLoggedIn)
	e.POST("/addUser", addUserHandler, middleware.EnsureLoggedIn, middleware.EnsureAdmin)
	e.POST("/removeUser", removeUserHandler, middleware.EnsureLoggedIn, middleware.EnsureAdmin)
	e.POST("/editUser", editUserHandler, middleware.EnsureLoggedIn, middleware.EnsureAdmin)
	e.POST("/readMessageHistory", readMessagesByPhoneNumberHandler, middleware.EnsureLoggedIn)
	e.POST("/signin", signinHandler)
	e.GET("/logout", logoutHandler, middleware.EnsureLoggedIn)
	e.GET("/settings", settingsHandler, middleware.EnsureLoggedIn, middleware.EnsureAdmin)
	e.POST("/sms", smsHandler)
	e.POST("/voice", voiceHandler)
	//e.POST("/welcome", welcomeHandler)
	//e.POST("/hold", holdHandler)
	e.POST("/connectAgent", connectAgentHandler)
	e.GET("/holdMusic", holdMusicHandler)
	e.POST("/dial", dialHandler, middleware.EnsureLoggedIn)
}
