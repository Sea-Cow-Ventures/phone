package main

import (
	"context"
	"net/http"

	"go.uber.org/zap"
	"golang.ngrok.com/ngrok"
	ngrokConfig "golang.ngrok.com/ngrok/config"
)

func startWebserverInNgrokTunnel() {
	tun, err := ngrok.Listen(context.Background(),
		ngrokConfig.HTTPEndpoint(
			ngrokConfig.WithDomain(cnf.UrlBasePath),
		),
		ngrok.WithAuthtoken(cnf.NgrokToken),
	)

	if err != nil {
		logger.Fatal("Failed to start ngrok tunnel")
	}
	logger.Info("Started ngrok tunnel")

	server := &http.Server{
		Handler: e,
	}

	// Start the server using the ngrok listener
	err = server.Serve(tun)
	if err != nil {
		logger.Fatal("Failed to start webserver with ngrok", zap.Error(err))
	}
	logger.Info("Started webserver")
}
