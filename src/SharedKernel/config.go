package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var C *Configuration

type Configuration struct {
	NatsConfig     *NatsConfig     `mapstructure:"nats"`
	DatabaseConfig *DatabaseConfig `mapstructure:"database"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

func (db *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.DBName, db.SSLMode)
}

type NatsConfig struct {
	URL            string        `mapstructure:"url"`
	MaxReconnects  int           `mapstructure:"max_reconnects"`
	ReconnectWait  time.Duration `mapstructure:"reconnect_wait"`
	ConnectTimeout time.Duration `mapstructure:"connect_timeout"`

	UserCredentials string `mapstructure:"user_credentials"`

	MaxPendingBytes int64 `mapstructure:"max_pending_bytes"`
	MaxPendingMsgs  int   `mapstructure:"max_pending_msgs"`
}

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	viper.AutomaticEnv()

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found. Using defaults and environment variables")
		} else {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Configuration
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("unable to decode config into struct: %w", err)
	}

	C = &config
	return nil
}

func setDefaults() {
	// NATS defaults
	viper.SetDefault("nats.url", "nats://localhost:6752")
	viper.SetDefault("nats.max_reconnects", 5)
	viper.SetDefault("nats.reconnect_wait", "2s")
	viper.SetDefault("nats.connect_timeout", "5s")
	viper.SetDefault("nats.max_pending_bytes", 67108864) // 64MB
	viper.SetDefault("nats.max_pending_msgs", 65536)     // 64k messages

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "admin")
	viper.SetDefault("database.dbname", "app_db")
	viper.SetDefault("database.sslmode", "disable")
}
