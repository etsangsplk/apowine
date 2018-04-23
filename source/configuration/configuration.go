package configuration

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	flag "github.com/spf13/pflag"
)

// Configuration stuct is used to populate the various fields used by apowine
type Configuration struct {
	ServerIP      string
	ServerPort    string
	ClientAddress string

	MongoUsername       string
	MongoPassword       string
	MongoDatabaseName   string
	MongoCollectionName string
	MongoURL            string
	DBSkipTLS           bool

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
	flag.String("ServerIP", "", "Server IP [Default: http://localhost]")
	flag.String("ServerPort", "", "Server Port [Default: 3000]")
	flag.String("ClientAddress", "", "Server Address [Default: 3005]")
	flag.String("LogLevel", "", "Log level. Default to info (trace//debug//info//warn//error//fatal)")
	flag.String("LogFormat", "", "Log Format. Default to human")

	flag.String("MongoUsername", "", "Username of the database [default: ]")
	flag.String("MongoPassword", "", "Password of the database [default: ]")
	flag.String("MongoDatabaseName", "", "Name of the database [default: drinksdb]")
	flag.String("MongoCollectionName", "", "Name of the collection in database [default: drinkscollection]")
	flag.String("MongoURL", "", "URI to connect to DB [default: 127.0.0.1:27017]")
	flag.Bool("DBSkipTLS", true, "Is valid TLS required for the DB server ? [default: true]")

	// Setting up default configuration
	viper.SetDefault("ServerIP", "http://localhost")
	viper.SetDefault("ServerPort", ":3000")
	viper.SetDefault("ClientAddress", ":3005")
	viper.SetDefault("LogLevel", "info")
	viper.SetDefault("LogFormat", "human")

	viper.SetDefault("MongoUsername", "")
	viper.SetDefault("MongoPassword", "")
	viper.SetDefault("MongoDatabaseName", "drinksdb")
	viper.SetDefault("MongoCollectionName", "drinkscollection")
	viper.SetDefault("MongoURL", "127.0.0.1:27017")
	viper.SetDefault("DBSkipTLS", true)

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
