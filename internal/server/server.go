package server

import (
	"aidan/phone/pkg/util"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"aidan/phone/internal/config"
	"aidan/phone/internal/log"
	customMiddleware "aidan/phone/internal/middleware"
	"aidan/phone/internal/models"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.ngrok.com/ngrok"
	ngrokConfig "golang.ngrok.com/ngrok/config"

	"go.uber.org/zap"
)

var (
	cnf    config.AppConfig
	logger *zap.Logger
)

func init() {
	cnf = config.GetConfig()
	logger = log.GetLogger()
}

func Start() {
	e := echo.New()

	if cnf.Env == "dev" {
		go createNgrokTunnel(e)
	} else {
		serveWithTls(e)
	}

	loadTemplates(e)
	go startLiveReloadWatcher(e)

	e.HideBanner = true
	e.HidePort = true

	e.Use(customMiddleware.Recover())
	e.Use(customMiddleware.Log())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))
	e.Validator = &models.Validator{Validator: validator.New()}

	RegisterRoutes(e)
}

func serveWithTls(e *echo.Echo) {
	// Build certificate and key file paths from the 'crt' folder.
	workingDir, err := util.GetWorkingDir()
	if err != nil {
		logger.Fatal("Failed to get working directory", zap.Error(err))
	}
	certFile := filepath.Join(workingDir, "crt/cert.pem")
	keyFile := filepath.Join(workingDir, "crt/key.pem")
	addr := ":443" // Listening on port 443 for HTTPS

	// Create the HTTP server with the Echo instance as the handler.
	server := &http.Server{
		Addr:    addr,
		Handler: e,
	}

	// Start the HTTPS server in a separate goroutine.
	go func() {
		if err := server.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start webserver with TLS", zap.Error(err))
		}
	}()

	logger.Info(fmt.Sprintf("Started TLS server on %s", addr))
}

func createNgrokTunnel(e *echo.Echo) {
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
}

func loadTemplates(e *echo.Echo) {
	parsedTemplates, err := template.New("").Funcs(template.FuncMap{
		"toJSON": util.ToJSON,
	}).ParseGlob(cnf.WebDir + "/templates/*.html")
	if err != nil {
		logger.Fatal("Failed to parse templates", zap.Error(err))
	}

	e.Renderer = &models.Template{
		Templates: parsedTemplates,
	}
}

func startLiveReloadWatcher(e *echo.Echo) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal("Failed to create fsnotify watcher", zap.Error(err))
	}
	defer watcher.Close()

	err = filepath.Walk(cnf.WebDir+"/templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		logger.Fatal("Failed to add template directories to watcher", zap.Error(err))
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) != 0 {
				logger.Info("Template change detected", zap.String("file", event.Name))
				loadTemplates(e)
			}
		case err := <-watcher.Errors:
			logger.Error("Error in fsnotify watcher", zap.Error(err))
		}
	}
}
