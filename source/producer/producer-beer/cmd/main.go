package main

import (
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/aporeto-inc/apowine/source/producer/producer-beer"
	"github.com/aporeto-inc/apowine/source/producer/producer-beer/configuration"
	"github.com/aporeto-inc/apowine/source/version"
)

func banner(version, revision string) {
	fmt.Printf(`


	  	  ___  ______ _____  _    _ _____ _   _  _____
		 / _ \ | ___ \  _  || |  | |_   _| \ | ||  ___|
		/ /_\ \| |_/ / | | || |  | | | | |  \| || |__
		|  _  ||  __/| | | || |/\| | | | | .\  ||  __|
		| | | || |   \ \_/ /\  /\  /_| |_| |\  || |___
		\_| |_/\_|    \___/  \/  \/ \___/\_| \_/\____/
    PRODUCER - BEER
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

	zap.L().Debug("Config used", zap.Any("Config", cfg.ServerURI))

	err = producerbeer.PushBeersToDB(cfg.ServerURI)
	if err != nil {
		log.Fatal("error adding beers to database", err)
	}

	zap.L().Info("Pushing beers to DB. Exiting in 5 seconds")

	<-time.After(time.Second * 5)

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
