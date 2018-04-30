package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/aporeto-inc/apowine/source/frontend-ui/client"
	"github.com/aporeto-inc/apowine/source/frontend-ui/client/internal/auth"
	"github.com/aporeto-inc/apowine/source/frontend-ui/client/internal/credential"
	"github.com/aporeto-inc/apowine/source/frontend-ui/configuration"
	"github.com/aporeto-inc/apowine/source/version"
	"github.com/gorilla/mux"
)

func banner(version, revision string) {
	fmt.Printf(`


	  	  ___  ______ _____  _    _ _____ _   _  _____
		 / _ \ | ___ \  _  || |  | |_   _| \ | ||  ___|
		/ /_\ \| |_/ / | | || |  | | | | |  \| || |__
		|  _  ||  __/| | | || |/\| | | | | .\  ||  __|
		| | | || |   \ \_/ /\  /\  /_| |_| |\  || |___
		\_| |_/\_|    \___/  \/  \/ \___/\_| \_/\____/
    CLIENT
_______________________________________________________________
             %s - %s
                                                 🚀  by Aporeto
`, version, revision)
}

func main() {
	banner(version.VERSION, version.REVISION)

	cfg, err := configuration.LoadConfiguration()
	if err != nil {
		log.Fatal("error parsing configuration", err)
	}

	err = setLogs(cfg.LogFormat, cfg.LogLevel)
	if err != nil {
		log.Fatalf("Error setting up logs: %s", err)
	}

	r := mux.NewRouter()

	zap.L().Debug("Config used", zap.Any("Config", cfg))

	handler := cors.Default().Handler(r)

	options := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler = options.Handler(handler)

	googleCreds := credential.NewGoogleCreds(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURI, cfg.GoogleRefreshToken)
	githubCreds := credential.NewGithubCreds(cfg.GithubClientID, cfg.GithubClientSecret, cfg.GithubRedirectURI)

	authHandler := auth.NewAuth(googleCreds, githubCreds)
	r.HandleFunc("/login", authHandler.Login)
	r.HandleFunc("/oauth2/github/callback", authHandler.GithubCallbackHandler).Methods(http.MethodGet)
	r.HandleFunc("/oauth2/google/callback", authHandler.GoogleCallbackHandler).Methods(http.MethodGet)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("/Users/sibi/apomux/workspace/code/go/src/github.com/aporeto-inc/apowine/source/frontend-ui/templates"))))

	clientHandler := client.NewClient(cfg.ServerAddress, authHandler)
	r.HandleFunc("/", client.GenerateLoginPage)
	r.HandleFunc("/home", clientHandler.GenerateClientPage)
	r.HandleFunc("/drink", clientHandler.GenerateDrinkManipulator)
	r.HandleFunc("/random", clientHandler.GenerateRandomDrinkManipulator)

	go func() {
		if err := http.ListenAndServe(cfg.ClientAddress, handler); err != nil {
			log.Fatal("error starting server", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	zap.L().Info("Everything started. Waiting for Stop signal")
	// Waiting for a Sig
	<-c

}

// setLogs setups Zap to log at the specified log level and format
func setLogs(logFormat, logLevel string) error {
	var zapConfig zap.Config

	switch logFormat {
	case "json":
		zapConfig = zap.NewProductionConfig()
		zapConfig.DisableStacktrace = true
	default:
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.DisableStacktrace = true
		zapConfig.DisableCaller = true
		zapConfig.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {}
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Set the logger
	switch logLevel {
	case "trace":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)
	return nil
}
