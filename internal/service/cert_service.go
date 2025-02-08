package service

import (
	"aidan/phone/internal/models"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/log"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/go-acme/lego/v4/registration"
	"go.uber.org/zap"
)

// LegoLogger implements the lego logger interface using zap
type LegoLogger struct {
	logger *zap.Logger
}

func (l *LegoLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(fmt.Sprint(args...))
}

func (l *LegoLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, args...))
}

func (l *LegoLogger) Fatalln(args ...interface{}) {
	l.logger.Fatal(fmt.Sprint(args...))
}

func (l *LegoLogger) Print(args ...interface{}) {
	l.logger.Info(fmt.Sprint(args...))
}

func (l *LegoLogger) Printf(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *LegoLogger) Println(args ...interface{}) {
	l.logger.Info(fmt.Sprint(args...))
}

func GenerateCert() error {
	logger := zap.L()
	logger.Info("Starting certificate generation")

	// Set up custom logger for lego
	log.Logger = &LegoLogger{logger: logger}

	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		logger.Error("Failed to generate private key", zap.Error(err))
		return err
	}

	user := &models.CertUser{
		Email:      cnf.CertEmail,
		PrivateKey: privateKey,
	}

	certConfig := lego.NewConfig(user)

	certConfig.CADirURL = lego.LEDirectoryProduction
	certConfig.Certificate.KeyType = certcrypto.EC384

	client, err := lego.NewClient(certConfig)
	if err != nil {
		return err
	}

	// Register user
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return fmt.Errorf("failed to register account: %w", err)
	}

	user.Registration = reg

	// Configure Cloudflare provider
	config := cloudflare.NewDefaultConfig()
	config.AuthEmail = cnf.CloudflareEmail
	config.AuthKey = cnf.CloudflareAPIKey
	//config.ZoneID = cnf.CloudflareZoneID

	provider, err := cloudflare.NewDNSProviderConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create Cloudflare provider: %w", err)
	}

	err = client.Challenge.SetDNS01Provider(provider)
	if err != nil {
		return fmt.Errorf("failed to set DNS provider: %w", err)
	}

	request := certificate.ObtainRequest{
		Domains: []string{cnf.UrlBasePath},
		Bundle:  true,
	}

	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return fmt.Errorf("failed to obtain certificate: %w", err)
	}

	err = os.MkdirAll("crt", 0755)
	if err != nil {
		return fmt.Errorf("failed to create crt directory: %w", err)
	}

	err = os.WriteFile("crt/cert.pem", certificates.Certificate, 0644)
	if err != nil {
		return fmt.Errorf("failed to write certificate: %w", err)
	}

	err = os.WriteFile("crt/key.pem", certificates.PrivateKey, 0600)
	if err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	err = os.WriteFile("crt/chain.pem", certificates.IssuerCertificate, 0644)
	if err != nil {
		return fmt.Errorf("failed to write certificate chain: %w", err)
	}

	return nil
}
