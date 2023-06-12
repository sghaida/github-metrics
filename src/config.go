package src

import (
	"github.com/spf13/viper"
	"log"
	"sync"
)

var (
	doOnce sync.Once
	config *Config
)

type Repo struct {
	Frontend []string `mapstructure:"frontend"`
	Backend  []string `mapstructure:"backend"`
}

type Team struct {
	Frontend []string `mapstructure:"frontend"`
	Backend  []string `mapstructure:"backend"`
}

type Config struct {
	Token string `mapstructure:"token" yaml:"token"`
	Org   string `mapstructure:"org" yaml:"org"`
	Repos Repo   `mapstructure:"repos" yaml:"repos"`
	Teams Team   `mapstructure:"teams" yaml:"teams"`
}

func GetConfig() *Config {
	doOnce.Do(func() {

		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("..")
		viper.AddConfigPath(".")

		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("Failed to read the config file: %v", err)
		}

		// Unmarshal the configuration into a struct
		err = viper.Unmarshal(&config)
		if err != nil {
			log.Fatalf("Failed to unmarshal config: %v", err)
		}
	})
	return config
}
