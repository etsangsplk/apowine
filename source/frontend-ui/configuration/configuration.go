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

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURI  string
	GoogleRefreshToken string

	GithubClientID     string
	GithubClientSecret string
	GithubRedirectURI  string

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
	flag.String("ClientAddress", "", "Server Address [Default: 3005]")
	flag.String("LogLevel", "", "Log level. Default to info (trace//debug//info//warn//error//fatal)")
	flag.String("LogFormat", "", "Log Format. Default to human")

	flag.String("GoogleClientID", "", "Google ClientID of the application [Default: 450167263420-ok2arm93kod5jdcecu5nkn93lmgjjjhh.apps.googleusercontent.com]")
	flag.String("GoogleClientSecret", "", "Google ClientSecret of the application [Default: LwttSwh51R1BKueI0w3n96dF]")
	flag.String("GoogleRedirectURI", "", "Google RedirectURI once authenticated [Default: http://localhost:3005/oauth2/google/callback]")
	flag.String("GoogleRefreshToken", "", "Google GoogleRefreshToken once authenticated [Default: 1/wEFTJhepVP8GHOPw1IJOa0P7F0q_OUn1pbLEpmjqGA6lP4fkEFWqujqzIA88L6mn]")

	flag.String("GithubClientID", "", "Github ClientID of the application [Default: 560f2688c98130ea6234]")
	flag.String("GithubClientSecret", "", "Github ClientSecret of the application [Default: 1c914abf80c3edd5b93cff3e47ab45afe96928bf]")
	flag.String("GithubRedirectURI", "", "Github RedirectURI once authenticated [Default: http://localhost:3000/oauth2/github/callback]")

	// Setting up default configuration
	viper.SetDefault("ServerAddress", "http://localhost:3000")
	viper.SetDefault("ClientAddress", ":3005")
	viper.SetDefault("LogLevel", "info")
	viper.SetDefault("LogFormat", "human")

	viper.SetDefault("GoogleClientID", "450167263420-ok2arm93kod5jdcecu5nkn93lmgjjjhh.apps.googleusercontent.com")
	viper.SetDefault("GoogleClientSecret", "LwttSwh51R1BKueI0w3n96dF")
	viper.SetDefault("GoogleRedirectURI", "http://localhost:3005/oauth2/google/callback")
	viper.SetDefault("GoogleRefreshToken", "1/wEFTJhepVP8GHOPw1IJOa0P7F0q_OUn1pbLEpmjqGA6lP4fkEFWqujqzIA88L6mn")

	viper.SetDefault("GithubClientID", "560f2688c98130ea6234")
	viper.SetDefault("GithubClientSecret", "1c914abf80c3edd5b93cff3e47ab45afe96928bf")
	viper.SetDefault("GithubRedirectURI", "http://localhost:3000/oauth2/github/callback")

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
