package config

import "os"

var config *Config

const (
	SERVICENAME = "phonebook"
	ENVIRONMENT = "ENVIRONMENT"
	HTTPPORT    = "HTTP_PORT"
	HTTPHOST    = "HTTP_HOST"
)

// Config phonebook application configuration
type Config struct {
	HTTPAddress        string
	Enviroment         string
	DBType             string
	DBConnectionString string
	MaxOpenCon         string
	MaxIdleCon         string
}

var defaultConfig = &Config{
	HTTPAddress:        getEnvOrDefault(HTTPHOST, "") + ":" + getEnvOrDefault(HTTPPORT, "9080"),
	Enviroment:         "development",
	DBType:             "postgres",
	DBConnectionString: "host=localhost port=5432 database=ledger user=root password=root sslmode=disable",
	MaxOpenCon:         "15",
	MaxIdleCon:         "10",
}

func getEnvOrDefault(env string, defaultVal string) string {
	e := os.Getenv(env)
	if e == "" {
		return defaultVal
	}
	return e
}

func Get() (*Config, error) {
	// return defaultConfig
	if config != nil {
		return config, nil
	}

	return defaultConfig, nil
}
