package config

import (
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
	WebDir          string   `json:"webDir"`
	MailServer      string   `json:"MailServer"`
	EmailRecipients []string `json:"EmailRecipients"`

	//CertFile         string   `json:"CertFile"`
	//CertKey          string   `json:"CertKey"`
	DBServer         string `json:"dbServer"`
	DBUser           string `json:"dbUser"`
	DBPass           string `json:"dbPass"`
	DBSchema         string `json:"dbSchema"`
	CertEmail        string `json:"certEmail"`
	CloudflareEmail  string `json:"cloudflareEmail"`
	CloudflareAPIKey string `json:"cloudflareApiKey"`
	CloudflareZoneID string `json:"cloudflareZoneID"`
}

var cnf AppConfig

func init() {
	configFile, err := os.Open("config.json")
	if err != nil {
		panic(fmt.Errorf("failed to open json config file: %w", err))
	}
	defer configFile.Close()

	data, err := io.ReadAll(configFile)
	if err != nil {
		panic(fmt.Errorf("failed to read json config file: %w", err))
	}

	var config AppConfig

	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal json config file: %w", err))
	}

	cnf = config
}

func GetConfig() AppConfig {
	return cnf
}
