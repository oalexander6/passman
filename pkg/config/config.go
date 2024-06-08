package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Host   string `mapstructure:"HOST" json:"HOST"`
	Debug  bool   `mapstructure:"DEBUG" json:"DEBUG"`
	EncIV  string `mapstructure:"ENC_IV" json:"ENC_IV"`
	EncKey string `mapstructure:"ENC_KEY" json:"ENC_KEY"`
}

var (
	activeConfig *Config
	lock         = &sync.Mutex{}
)

func GetConfig() *Config {
	// extra check here to avoid using the (very expensive) lock whenever possible
	if activeConfig == nil {
		lock.Lock()
		defer lock.Unlock()

		if activeConfig == nil {
			activeConfig = new()
		}
	}

	return activeConfig
}

func new() *Config {
	c := &Config{}

	v := viper.New()

	v.SetEnvPrefix("PASSMAN")

	v.SetDefault("HOST", "localhost:8080")
	v.SetDefault("DEBUG", false)

	v.SetConfigName("passman-config")
	v.AddConfigPath("/etc/passman")
	v.AddConfigPath("$HOME/.passman")

	projectDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Somehow failed to get the working directory")
	}
	v.AddConfigPath(projectDir)

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Fatal("Failed to load the configuration: ", err)
	}

	if err := v.Unmarshal(c); err != nil {
		log.Fatal("Failed to parse configuration: ", err)
	}

	if err = c.validate(); err != nil {
		log.Fatal("Invalid configuration: ", err)
	}

	log.Default().Println("Conductor configuration initialized")

	vals, _ := json.MarshalIndent(c, "", "\t")
	log.Default().Println(string(vals))

	return c
}

func (c *Config) validate() error {
	if len(c.EncIV) != 16 {
		return errors.New("must provide a 16 byte ENC_IV")
	}

	if len(c.EncKey) != 32 {
		return errors.New("must provide a 32 byte ENC_KEY")
	}

	return nil
}
