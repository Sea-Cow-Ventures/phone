package middleware

import (
	"aidan/phone/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func Recover() echo.MiddlewareFunc {
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

					return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error(), Success: false})
				}
			}
			return next(c)
		}
	}
}
