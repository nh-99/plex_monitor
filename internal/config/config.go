package config

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Config is the global configuration for the application. It can be loaded from a file or environment variables.
// It is a singleton and should be accessed via the GetConfig() function. It is initialized in the main function.
// It is a struct so that it can be extended in the future.
type Config struct {
	// Database is the configuration for the database
	Database DatabaseConfig `yaml:"database"`

	// Server is the configuration for the server
	Server ServerConfig `yaml:"server"`

	// Discord is the configuration for the Discord bot
	Discord DiscordConfig `yaml:"discord"`

	// Logging is the configuration for the logging
	Logging AppLogging `yaml:"logging"`

	// Environment is the environment that the application is running in
	Environment Environment `yaml:"environment"`

	// IsProduction is true if the application is running in production mode
	IsProduction bool `yaml:"isProduction"`
}

// DatabaseConfig is the configuration for the database
type DatabaseConfig struct {
	// ConnectionString is the connection string for the database
	ConnectionString string `yaml:"connectionString"`
	// Name is the name of the database
	Name string `yaml:"name"`
}

// ServerConfig is the configuration for the server
type ServerConfig struct {
	// Port is the port that the server will listen on
	Port string `yaml:"port"`
}

// DiscordConfig is the configuration for the Discord bot
type DiscordConfig struct {
	// Token is the token that is used to authenticate with Discord
	Token string `yaml:"token"`
}

// AppLogging is the configuration for the application logging
type AppLogging struct {
	// Level is the log level that the application should log at
	Level string `yaml:"level" default:"info"`
	// ReportCaller is true if the caller should be reported in the logs
	ReportCaller bool `yaml:"reportCaller" default:"true"`
	// Format is the format that the logs should be in
	Format string `yaml:"format" default:"json"`
}

// GetLogrusLevel returns the logrus level that the application should log at
func (al *AppLogging) GetLogrusLevel() logrus.Level {
	level, err := logrus.ParseLevel(al.Level)
	if err != nil {
		panic(fmt.Errorf("invalid log level %s: %w", al.Level, err))
	}
	return level
}

// Environment is the environment that the application is running in
type Environment string

const (
	// EnvironmentDev is the development environment
	EnvironmentDev Environment = "dev"
	// EnvironmentProd is the production environment
	EnvironmentProd Environment = "prod"
)

func init() {
	if config == (Config{}) {
		populateConfig()
	}
}

var config Config

func getEnvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("Environment variable " + key + " is not set")
	}
	return value
}

func populateConfig() {
	// Check if the database URL is set. If it is, load all config from the environment variables. Otherwise, load from
	// the config file.
	if os.Getenv("DATABASE_URL") != "" {
		SetConfig(Config{
			Database: DatabaseConfig{
				ConnectionString: getEnvOrPanic("DATABASE_URL"),
				Name:             getEnvOrPanic("DATABASE_NAME"),
			},
			Server: ServerConfig{
				Port: getEnvOrDefault("PORT", "8080"),
			},
			Discord: DiscordConfig{
				Token: getEnvOrPanic("DISCORD_BOT_TOKEN"),
			},
			Environment:  Environment(getEnvOrDefault("ENVIRONMENT", "dev")),
			IsProduction: Environment(getEnvOrDefault("ENVIRONMENT", "dev")) == EnvironmentProd,
		})
	} else {
		// Load config from yaml file
		yamlFile, err := os.ReadFile("config.yaml")
		if err != nil {
			path, _ := os.Getwd()
			panic(fmt.Errorf("error reading config file from path=%s: %w", path, err))
		}

		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			panic(err)
		}
	}
}

// GetConfig returns the global configuration for the application
func GetConfig() Config {
	return config
}

// SetConfig sets the global configuration for the application
func SetConfig(c Config) {
	config = c
}
