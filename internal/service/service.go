package service

import (
	"aidan/phone/internal/config"
	"aidan/phone/internal/log"
	"fmt"

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
	var err error
	cnf = config.GetConfig()
	temp := cnf.TwilioUser
	fmt.Println("temp", temp)
	logger = log.GetLogger()
	t, err = ConnectTwilio()
	if err != nil {
		//logger.Fatal("Unable to connect to twilio", zap.Error(err))
	}
}

func ConnectTwilio() (*twilio.RestClient, error) {
	t := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cnf.TwilioUser,
		Password: cnf.TwilioPass,
	})

	fmt.Println("cnf", cnf.TwilioUser, cnf.TwilioPass)
	//logger.Info("cnf", zap.Any("cnf", cnf))
	createParams := &api.CreateIncomingPhoneNumberParams{}
	createParams.SetPhoneNumber(cnf.PhoneNumber)
	createParams.SetSmsUrl(cnf.UrlBasePath + "/sms")
	createParams.SetVoiceUrl(cnf.UrlBasePath + "/voice")

	resp, err := t.Api.CreateIncomingPhoneNumber(createParams)
	if err != nil {
		return nil, fmt.Errorf("Unable to create twilio phone number", err)
	}

	logger.Info("Connected to twilio", zap.String("phoneNumber", cnf.PhoneNumber), zap.Any("twilioResponse", resp))
	return t, nil
}
