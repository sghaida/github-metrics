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
	Type  TeamType `mapstructure:"type"`
	Names []string `mapstructure:"names"`
}

type Team struct {
	Name    TeamType `mapstructure:"name"`
	Members []string `mapstructure:"members"`
}

type Config struct {
	Tokens []string          `mapstructure:"tokens" yaml:"token"`
	Org    string            `mapstructure:"org" yaml:"org"`
	Repos  []Repo            `mapstructure:"repos" yaml:"repos"`
	Teams  map[string][]Team `mapstructure:"teams" yaml:"teams"`
}

func GetConfig() *Config {
	doOnce.Do(func() {

		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("..")
		viper.AddConfigPath("/app")

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
