package main

import (
	"fmt"
	"time"

	"aidan/phone/internal/config"
	"aidan/phone/internal/log"
	"aidan/phone/internal/service"
	"aidan/phone/pkg/util"

	"go.uber.org/zap"
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

var (
	cnf    config.AppConfig
	logger *zap.SugaredLogger
	//serviceName string = "TaskScheduler"
	//fromEmail   string = "TaskScheduler@fweco.net"
	workingDir string
	//rootCA      []byte
	//tasks       []Task
)

var inQueue int

func init() {
	var err error
	workingDir, err = util.GetWorkingDir()
	if err != nil {
		panic(fmt.Errorf("unable to get working dir: %w", err))
	}

	log.Init(workingDir, "phone.log")

	logger = log.GetLogger()

	logger.Info("Init")

	err = config.LoadConfig(workingDir + "/config.json")
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	cnf = config.GetConfig()

	//server.Start()

	//initMysql()

	/*agents, err := agent.ReadAgents()
	if err != nil {
		logger.Fatal("Failed to read agents", zap.Error(err))
	}

	if len(agents) < 1 {
		admin, err := agent.CreateDefaultAdmin()
		if err != nil {
			logger.Fatal("Failed to create default admin", zap.Error(err))
		}
		logger.Info("Created default admin", zap.Any("details", admin))
	}*/

}

func main() {
	logger.Info("Ready")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		go func() {
			if err := service.GetAccountCallHistory(); err != nil {
				logger.Errorf("Error reading call history: %v", err)
			}
		}()

		go func() {
			if err := service.GetAccountMessageHistory(); err != nil {
				logger.Errorf("Error reading message history: %v", err)
			}
		}()
	}
}
