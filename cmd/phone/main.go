package main

import (
	"time"

	"aidan/phone/internal/log"
	"aidan/phone/internal/server"
	"aidan/phone/internal/service"

	"go.uber.org/zap"
)

/*
	type VoiceMessage struct {
		Event          string `json:"event"`
		SequenceNumber string `json:"sequenceNumber"`
		StreamSid      string `json:"streamSid"`
	}

	type MediaFormat struct {
		Track     string `json:"track"`
		Chunk     string `json:"chunk"`
		Timestamp string `json:"timestamp"`
		Payload   string `json:"payload"`
	}
*/
func main() {
	logger := log.GetLogger()

	err := service.GenerateCert()
	if err != nil {
		logger.Fatal("Failed getting lets encrypt certificate", zap.Error(err))
	}

	server.Start()

	logger.Info("Ready")

	if adminCount, err := service.GetAdminCount(); err != nil || adminCount == 0 {
		logger.Info("Creating default admin agent")
		admin, err := service.CreateDefaultAdmin()
		if err != nil {
			logger.Fatal("Error creating default admin agent", zap.Error(err))
		}
		logger.Info("Default admin agent created", zap.Any("admin", admin))
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		go func() {
			if err := service.GetAccountCallHistory(); err != nil {
				logger.Sugar().Errorf("Error reading call history: %v", err)
			}
		}()

		go func() {
			if err := service.GetAccountMessageHistory(); err != nil {
				logger.Sugar().Errorf("Error reading message history: %v", err)
			}
		}()
	}
}
