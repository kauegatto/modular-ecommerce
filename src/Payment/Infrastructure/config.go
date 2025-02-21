package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var C *Configuration

type Configuration struct {
	ERedeConfig *ERedeConfig `mapstructure:"eRede"`
}

type ERedeConfig struct {
	PV            string `mapstructure:"pv"`
	Authorization string `mapstructure:"authorization"`
	BaseURL       string `mapstructure:"baseURL"`
	Timeout       int    `mapstructure:"timeout"`
}

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	viper.AutomaticEnv()

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
