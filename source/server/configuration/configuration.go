package configuration

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	flag "github.com/spf13/pflag"
)

// Configuration stuct is used to populate the various fields used by apowine
type Configuration struct {
	ServerPort string

	UseHealth  bool
	HealthPort string

	MakeNewConnection   bool
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

	flag.String("ServerPort", ":3000", "Server Port ")
	flag.String("LogLevel", "info", "Log level. (trace//debug//info//warn//error//fatal)")
	flag.String("LogFormat", "human", "Log Format. ")

	flag.Bool("MakeNewConnection", true, "To create new session every request ")
	flag.String("MongoDatabaseName", "drinksdb", "Name of the database ")
	flag.String("MongoCollectionName", "drinkscollection", "Name of the collection in database ")
	flag.String("MongoURL", "127.0.0.1:27017", "URI to connect to DB ")
	flag.Bool("DBSkipTLS", true, "Is valid TLS required for the DB server ? ")

	flag.Bool("UseHealth", false, "Use health ")
	flag.Int("HealthPort", 5000, "Health Port ")

	viper.SetDefault("AuthorizedEmail", "aliceaporeto@gmail.com")
	viper.SetDefault("AuthorizedGivenName", "Alice")
	viper.SetDefault("AuthorizedFamilyName", "Aporeto")

	viper.SetDefault("UnauthorizedEmail", "bobaporeto@gmail.com")
	viper.SetDefault("UnsuthorizedGivenName", "Bob")
	viper.SetDefault("UnauthorizedFamilyName", "Aporeto")

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
