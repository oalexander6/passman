package config

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

const (
	LOCAL_ENV = "LOCAL"
	DEV_ENV   = "DEV"
	STAGE_ENV = "STAGE"
	PROD_ENV  = "PROD"
)

type PostgresConfig struct {
	// Postgres host
	Host string `json:"HOST" validate:"required"`
	// Postgres database name
	DBName string `json:"DB_NAME" validate:"required"`
	// Postgres username
	User string `json:"USER" validate:"required"`
	// Postgres password
	Password string `json:"-" validate:"required"`
}

type EncryptionConfig struct {
	// Initialization vector for AES encryption
	EncIV string `json:"ENCRYPTION_IV" validate:"required"`
	// AES encryption secret key
	EncSecret string `json:"ENCRYPTION_SECERET" validate:"required"`
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
	// CSRF token secret key
	CSRFKey string `json:"-" validate:"required_if=EnableCSRFProtection true,len=36"`
	// use CSRF protections
	EnableCSRFProtection bool `json:"ENABLE_CSRF_PROTECTION" validate:"boolean"`
	// Postgres configuration
	PostgresOpts PostgresConfig `json:"POSTGRES" validate:"required"`
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
		CSRFKey:   secretVals["CSRF_KEY"],
		PostgresOpts: PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			DBName:   os.Getenv("POSTGRES_DB_NAME"),
			User:     secretVals["POSTGRES_USER"],
			Password: secretVals["POSTGRES_PASSWORD"],
		},
	}

	useCSRF, err := strconv.ParseBool(os.Getenv("ENABLE_CSRF_PROTECTION"))
	if err != nil {
		panic("Failed to parse value for ENABLE_CSRF_PROTECTION as a bool")
	}
	c.EnableCSRFProtection = useCSRF

	return c
}

func loadSecrets() (map[string]string, error) {
	loadedVals := make(map[string]string)

	secrets := []string{"SECRET_KEY", "POSTGRES_USER", "POSTGRES_PASSWORD", "CSRF_KEY", "ENCRYPTION_IV", "ENCRYPTION_SECRET"}

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
