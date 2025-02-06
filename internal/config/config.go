package config

import (
	"aidan/phone/pkg/util"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

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

var cnf AppConfig

func LoadConfig(path string) error {
	configFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open json config file: %w", err)
	}
	defer configFile.Close()

	data, err := io.ReadAll(configFile)
	if err != nil {
		return fmt.Errorf("failed to read json config file: %w", err)
	}

	var config AppConfig

	err = json.Unmarshal(data, config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json config file: %w", err)
	}

	cnf = config
	return nil
}

func GetConfig() AppConfig {
	if cnf.Env == "" && cnf.TwilioUser == "" {
		workingDir, err := util.GetWorkingDir()
		if err != nil {
			panic(fmt.Errorf("unable to get working dir: %w", err))
		}

		err = LoadConfig(workingDir + "/config.json")
		if err != nil {
			panic(fmt.Errorf("failed to load config: %w", err))
		}
	}
	return cnf
}
