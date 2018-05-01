package configuration

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	flag "github.com/spf13/pflag"
)

// Configuration stuct is used to populate the various fields used by apowine
type Configuration struct {
	ServerAddress string
	ClientAddress string

	MidgardTokenRealm    string
	MidgardTokenValidity string

	LogFormat string
	LogLevel  string
}

func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

// LoadConfiguration will load the configuration struct
func LoadConfiguration() (*Configuration, error) {
	flag.Usage = usage
	flag.String("ServerAddress", "", "Server IP [Default: http://localhost:3000]")
	flag.String("ClientAddress", "", "Server Address [Default: 43245]")
	flag.String("LogLevel", "", "Log level. Default to info (trace//debug//info//warn//error//fatal)")
	flag.String("LogFormat", "", "Log Format. Default to human")

	flag.String("MidgardTokenRealm", "", "Midgard realm [Default: Google]")
	flag.String("MidgardTokenValidity", "", "Midgard token validity [Default: 720h]")

	// Setting up default configuration
	viper.SetDefault("ServerAddress", "http://localhost:3000")
	viper.SetDefault("ClientAddress", ":43245")
	viper.SetDefault("LogLevel", "info")
	viper.SetDefault("LogFormat", "human")

	viper.SetDefault("MidgardTokenRealm", "Google")
	viper.SetDefault("MidgardTokenValidity", "720h")

	// Binding ENV variables
	// Each config will be of format TRIREME_XYZ as env variable, where XYZ
	// is the upper case config.
	viper.SetEnvPrefix("APOWINE")
	viper.AutomaticEnv()

	// Binding CLI flags.
	flag.Parse()
	viper.BindPFlags(flag.CommandLine)

	var config Configuration

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling: %s", err)
	}

	return &config, nil
}
