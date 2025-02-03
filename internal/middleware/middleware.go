package middleware

import (
	"aidan/phone/internal/log"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func init() {
	logger = log.GetLogger()
}
