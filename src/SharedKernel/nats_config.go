package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var NatsConfiguration *NatsConfig

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
	viper.AutomaticEnv()
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found. Using defaults and environment variables")
		} else {
			return fmt.Errorf("error reading config file: %s", err)
		}
	}

	config := &NatsConfig{}
	if err := viper.Unmarshal(config); err != nil {
		return fmt.Errorf("unable to decode config into struct: %s", err)
	}

	NatsConfiguration = config

	return nil
}

func setDefaults() {
	viper.SetDefault("url", "nats://localhost:6752")
	viper.SetDefault("max_reconnects", 5)
	viper.SetDefault("reconnect_wait", "2s")
	viper.SetDefault("connect_timeout", "5s")
	viper.SetDefault("max_pending_bytes", 67108864) // 64MB
	viper.SetDefault("max_pending_msgs", 65536)     // 64k messages
}
