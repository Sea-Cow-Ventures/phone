package handlers

import (
	"aidan/phone/internal/log"

	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger = log.GetLogger()
}
