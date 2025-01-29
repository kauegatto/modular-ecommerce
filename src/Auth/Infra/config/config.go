package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Port         int
	Host         string
	ReadTimeout  int
	WriteTimeout int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type LoggerConfig struct {
	Level string
	File  string
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found. Using defaults and environment variables")
		} else {
			return nil, fmt.Errorf("error reading config file: %s", err)
		}
	}

	config := &Config{}
	err := viper.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %s", err)
	}

	return config, nil
}

func setDefaults() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.readtimeout", 10)
	viper.SetDefault("server.writetimeout", 10)

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.dbname", "myapp")
	viper.SetDefault("database.sslmode", "disable")

	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.file", "app.log")
}
