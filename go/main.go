package main

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
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

type AppConfig struct {
	Env             string   `json:"env"`
	NgrokToken      string   `json:"ngrokToken"`
	TwilioUser      string   `json:"twilioUser"`
	TwilioPass      string   `json:"twilioPass"`
	UrlBasePath     string   `json:"urlBasePath"`
	ListenPort      string   `json:"listenPort"`
	PhoneNumber     string   `json:"phoneNumber"`
	HoldMusicPath   string   `json:"holdMusicPath"`
	MailServer      string   `json:"MailServer"`
	EmailRecipients []string `json:"EmailRecipients"`

	//CertFile         string   `json:"CertFile"`
	//CertKey          string   `json:"CertKey"`
	DBServer string `json:"dbServer"`
	DBUser   string `json:"dbUser"`
	DBPass   string `json:"dbPass"`
	DBSchema string `json:"dbSchema"`
}

var (
	cnf    AppConfig
	logger *zap.Logger
	e      *echo.Echo
	db     *sqlx.DB
	t      *twilio.RestClient
	//serviceName string = "TaskScheduler"
	//fromEmail   string = "TaskScheduler@fweco.net"
	workingDir string
	//rootCA      []byte
	//tasks       []Task
)

var inQueue int

func init() {
	var err error
	workingDir, err = getWorkingDir()
	if err != nil {
		panic(fmt.Errorf("unable to get working dir: %w", err))
	}

	logger = createLogger(workingDir, "go-phone.log")

	logger.Info("Init")

	err = loadConfig(workingDir+"/config.json", &cnf)
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	e = echo.New()

	if cnf.Env == "dev" {
		go startWebserverInNgrokTunnel()
	} else {
		//startWebserver()
	}

	initWebserver()

	initMysql()

	agents, err := readAgents()
	if err != nil {
		logger.Fatal("Failed to read agents", zap.Error(err))
	}

	if len(agents) < 1 {
		admin, err := createDefaultAdmin()
		if err != nil {
			logger.Fatal("Failed to create default admin", zap.Error(err))
		}
		logger.Info("Created default admin", zap.Any("details", admin))
	}

	t = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cnf.TwilioUser,
		Password: cnf.TwilioPass,
	})

	initTwilio()
}

func main() {
	logger.Info("Ready")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		go func() {
			if err := readAccountCallHistory(); err != nil {
				logger.Sugar().Errorf("Error reading call history: %v", err)
			}
		}()

		go func() {
			if err := readAccountMessageHistory(); err != nil {
				logger.Sugar().Errorf("Error reading message history: %v", err)
			}
		}()
	}
}

func initMysql() {
	cfg := mysql.Config{
		User:                 cnf.DBUser,
		Passwd:               cnf.DBPass,
		Net:                  "tcp",
		Addr:                 cnf.DBServer + ":3306",
		DBName:               cnf.DBSchema,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	mysql, err := sqlx.Open("mysql", cfg.FormatDSN())

	mysql.SetConnMaxLifetime(time.Minute * 3)
	mysql.SetMaxOpenConns(10)
	mysql.SetMaxIdleConns(10)

	if err != nil {
		logger.Fatal("Failed to connect to mysql db", zap.Error(err))
	}

	err = mysql.Ping()
	if err != nil {
		logger.Fatal("Failed to connect to mysql db", zap.Error(err))
	}

	db = mysql
	logger.Info("Connected to db")
}

func initTwilio() {
	createParams := &twilioApi.CreateIncomingPhoneNumberParams{}
	createParams.SetPhoneNumber(cnf.PhoneNumber)
	createParams.SetSmsUrl(cnf.UrlBasePath + "/sms")
	createParams.SetVoiceUrl(cnf.UrlBasePath + "/voice")

	resp, err := t.Api.CreateIncomingPhoneNumber(createParams)
	if err != nil {
		logger.Fatal("Unable to create twilio phone number", zap.Error(err))
	}

	logger.Info("Connected to twilio", zap.String("phoneNumber", cnf.PhoneNumber), zap.Any("twilioResponse", resp))
}
