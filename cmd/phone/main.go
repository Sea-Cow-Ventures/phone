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

	/*doesExist, err := service.DoesCertExist()
	if err != nil {
		logger.Fatal("Error checking if cert exists", zap.Error(err))
	}

	if !doesExist {
		logger.Info("Cert does not exist, generating...")
		err := service.GenerateCert()
		if err != nil {
			logger.Fatal("Error generating cert", zap.Error(err))
		}
	} else {
		isExpiringSoon, err := service.IsCertExpiringSoon()
		if err != nil {
			logger.Fatal("Error checking if cert is expiring soon", zap.Error(err))
		}
		if isExpiringSoon {
			logger.Info("Cert is expiring soon, renewing...")
			err := service.GenerateCert()
			if err != nil {
				logger.Fatal("Error renewing cert", zap.Error(err))
			}
		}
	}

	if adminCount, err := service.GetAdminCount(); err != nil || adminCount == 0 {
		logger.Info("Creating default admin agent")
		admin, err := service.CreateDefaultAdmin()
		if err != nil {
			logger.Fatal("Error creating default admin agent", zap.Error(err))
		}
		logger.Info("Default admin agent created", zap.Any("admin", admin))
	}*/

	webpushPublicKey, _, _ := service.GetWebpushKeys()
	if webpushPublicKey == "" {
		logger.Info("Generating webpush keys")
		_, _, err := service.GenerateWebpushKeys()
		if err != nil {
			logger.Fatal("Error generating webpush keys", zap.Error(err))
		}
	}

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
	/*var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10 * time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()*/
}
