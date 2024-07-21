package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Host              string `mapstructure:"HOST" json:"HOST"`
	Debug             bool   `mapstructure:"DEBUG" json:"DEBUG"`
	EncIV             string `mapstructure:"ENC_IV" json:"-"`
	EncKey            string `mapstructure:"ENC_KEY" json:"-"`
	SecretKey         string `mapstructure:"SECRET_KEY" json:"-"`
	UseCSRFTokens     bool   `mapstructure:"USE_CSRF_TOKENS" json:"USE_CSRF_TOKENS"`
	CSRFSecret        string `mapstructure:"CSRF_SECRET" json:"-"`
	StaticFilePath    string `mapstructure:"STATIC_FILE_PATH" json:"STATIC_FILE_PATH"`
	AvatarStoragePath string `mapstructure:"AVATAR_STORAGE_PATH" json:"AVATAR_STORAGE_PATH"`
	UseSSL            bool   `mapstructure:"USE_SSL" json:"USE_SSL"`
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

	v.SetDefault("HOST", "127.0.0.1:8080")
	v.SetDefault("DEBUG", false)
	v.SetDefault("USE_CSRF_TOKENS", true)
	v.SetDefault("AVATAR_STORAGE_PATH", "/tmp")
	v.SetDefault("USE_SSL", true)

	v.SetConfigName("passman-config")
	v.AddConfigPath("/etc/passman/config")
	v.AddConfigPath("$HOME/.passman")
	v.AddConfigPath("$PWD/config")

	projectDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Somehow failed to get the working directory")
	}
	v.SetDefault("STATIC_FILE_PATH", fmt.Sprintf("%s/dist/public/assets", projectDir))

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

	if len(c.SecretKey) < 36 {
		return errors.New("must provide a minimum 36 byte SECRET_KEY")
	}

	if c.UseCSRFTokens && len(c.CSRFSecret) < 36 {
		return errors.New("must provide a minimum 36 byte CSRF_SECRET if USE_CSRF_TOKENS is true")
	}

	return nil
}
