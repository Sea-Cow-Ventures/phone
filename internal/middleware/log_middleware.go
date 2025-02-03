package middleware

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Log() echo.MiddlewareFunc {
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

				logger.Errorw("Webserver error", fields)
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

				logger.Warnw("Webserver client error", fields)
			case n >= 300:
				logger.Infow("Webserver redirection", fields)
			default:
				logger.Infow("Webserver success", fields)
			}

			return nil
		}
	}
}
