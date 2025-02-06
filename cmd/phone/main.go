package main

import (
	"time"

	"aidan/phone/internal/log"
	"aidan/phone/internal/server"
	"aidan/phone/internal/service"
)

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

//var inQueue int

func main() {
	logger := log.GetLogger()

	server.Start()

	logger.Info("Ready")

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
