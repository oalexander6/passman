package config

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

const (
	LOCAL_ENV           = "LOCAL"
	DEV_ENV             = "DEV"
	STAGE_ENV           = "STAGE"
	PROD_ENV            = "PROD"
	STORE_TYPE_POSTGRES = "postgres"
	STORE_TYPE_SQLITE   = "sqlite"
)

type PostgresConfig struct {
	// Postgres connection URI
	URI string `json:"-" validate:"required"`
}

type EncryptionConfig struct {
	// Initialization vector for AES encryption
	EncIV string `json:"-" validate:"required,len=16"`
	// AES encryption secret key
	EncSecret string `json:"-" validate:"required,len=32"`
}

type Config struct {
	// LOCAL, DEV, STAGE, PROD
	Env string `json:"ENV" validate:"required,oneof=LOCAL DEV STAGE PROD"`
	// server listen port
	Port string `json:"PORT" validate:"required,numeric"`
	// current application version
	Version string `json:"VERSION" validate:"required"`
	// encryption key for sessions
	SecretKey string `json:"-" validate:"required"`
	// store type to use - postgres, sqlite
	StoreType string `json:"STORE_TYPE" validate:"required,oneof=postgres sqlite"`
	// Postgres configuration
	PostgresOpts PostgresConfig `json:"POSTGRES" validate:"required_if=StoreType postgres"`
	// Note encryption config
	Encryption EncryptionConfig `json:"ENCRYPTION" validate:"required"`
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file loaded, config will check existing env variables")
	}

	secretVals, err := loadSecrets()
	if err != nil {
		panic("Failed to load secret values")
	}

	c := &Config{
		Env:       strings.ToUpper(os.Getenv("ENV")),
		Port:      os.Getenv("PORT"),
		Version:   os.Getenv("VERSION"),
		SecretKey: secretVals["SECRET_KEY"],
		StoreType: os.Getenv("STORE_TYPE"),
		PostgresOpts: PostgresConfig{
			URI: os.Getenv("DB_URI"),
		},
		Encryption: EncryptionConfig{
			EncIV:     secretVals["ENCRYPTION_IV"],
			EncSecret: secretVals["ENCRYPTION_SECRET"],
		},
	}

	return c
}

func loadSecrets() (map[string]string, error) {
	loadedVals := make(map[string]string)

	secrets := []string{"SECRET_KEY", "DB_URI", "ENCRYPTION_IV", "ENCRYPTION_SECRET"}

	for _, baseEnvName := range secrets {
		// default to non-file variable if provided
		val := os.Getenv(baseEnvName)
		if val != "" {
			loadedVals[baseEnvName] = val
			continue
		}

		// if non-file version was not found, try the file version
		fileEnvVarName := baseEnvName + "_FILE"
		pathToLoad := os.Getenv(fileEnvVarName)

		if pathToLoad != "" {
			val, err := os.ReadFile(pathToLoad)
			if err != nil {
				return nil, err
			}
			loadedVals[baseEnvName] = string(val)
		}
	}

	return loadedVals, nil
}

func (c Config) Validate() error {
	var Validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

	if err := Validate.Struct(c); err != nil {
		return err
	}

	if !slices.Contains([]string{LOCAL_ENV, DEV_ENV, STAGE_ENV, PROD_ENV}, c.Env) {
		return fmt.Errorf("invalid env: %s", c.Env)
	}

	return nil
}
