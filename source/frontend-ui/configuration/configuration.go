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

	MidgardURL           string
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
	flag.String("ServerAddress", "http://localhost:3000", "Server IP ")
	flag.String("ClientAddress", ":3005", "Server Address ")
	flag.String("LogLevel", "info", "Log level. (trace//debug//info//warn//error//fatal)")
	flag.String("LogFormat", "human", "Log Format. ")

	flag.String("MidgardURL", "https://api.console.aporeto.com/issue", "URL of midgard server ")
	flag.String("MidgardTokenRealm", "Google", "Midgard realm ")
	flag.String("MidgardTokenValidity", "720h", "Midgard token validity ")

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
