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

			serverErr := next(c)
			if serverErr != nil {
				c.Error(serverErr)
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

			body, err := io.ReadAll(c.Request().Body)
			if err != nil {
				logger.Error("Failed to read request body", zap.Error(err))
				return err
			}
			c.Request().Body = io.NopCloser(bytes.NewReader(body))

			if len(body) > 0 {
				fields = append(fields, zap.String("request_body", string(body)))
			}

			n := res.Status
			switch {
			case n >= 500:
				fields = append(fields, zap.Error(serverErr))
				logger.Error("Webserver error", fields...)
			case n >= 400:
				logger.Warn("Webserver client error", fields...)
			case n >= 300:
				logger.Info("Webserver redirection", fields...)
			default:
				logger.Info("Webserver success", fields...)
			}

			return nil
		}
	}
}
