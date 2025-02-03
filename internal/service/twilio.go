package twilioWrapper

import (
	"aidan/phone/internal/config"
	"aidan/phone/internal/log"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"

	"go.uber.org/zap"
)

var (
	cnf    config.AppConfig
	t      *twilio.RestClient
	logger *zap.SugaredLogger
)

func init() {
	cnf = config.GetConfig()
	logger = log.GetLogger()
	t = Connect()
}

func Connect() *twilio.RestClient {
	t := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cnf.TwilioUser,
		Password: cnf.TwilioPass,
	})
	createParams := &api.CreateIncomingPhoneNumberParams{}
	createParams.SetPhoneNumber(cnf.PhoneNumber)
	createParams.SetSmsUrl(cnf.UrlBasePath + "/sms")
	createParams.SetVoiceUrl(cnf.UrlBasePath + "/voice")

	resp, err := t.Api.CreateIncomingPhoneNumber(createParams)
	if err != nil {
		logger.Fatal("Unable to create twilio phone number", zap.Error(err))
	}

	logger.Info("Connected to twilio", zap.String("phoneNumber", cnf.PhoneNumber), zap.Any("twilioResponse", resp))
	return t
}

func GetTwilioClient() *twilio.RestClient {
	return t
}
