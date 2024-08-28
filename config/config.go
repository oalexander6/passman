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
	LOCAL_ENV           = "LOCAL"
	DEV_ENV             = "DEV"
	STAGE_ENV           = "STAGE"
	PROD_ENV            = "PROD"
	STORE_TYPE_POSTGRES = "postgres"
	STORE_TYPE_SQLITE   = "sqlite"
)

type PostgresConfig struct {
	// Postgres connection URI
	URI string `json:"DB_URI" validate:"required"`
}

type SqliteConfig struct {
	// the file path to the database
	DBFile string `json:"DB_FILE" validate:"required"`
	// whether to delete all existing data on startup
	DeleteOnStartup bool `json:"DELETE_ON_STARTUP" validate:"required"`
}

type EncryptionConfig struct {
	// Initialization vector for AES encryption
	EncIV string `json:"ENCRYPTION_IV" validate:"required,len=16"`
	// AES encryption secret key
	EncSecret string `json:"ENCRYPTION_SECERET" validate:"required,len=32"`
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
	// store type to use - postgres, sqlite
	StoreType string `json:"STORE_TYPE" validate:"required,oneof=postgres sqlite"`
	// Postgres configuration
	PostgresOpts PostgresConfig `json:"POSTGRES" validate:"required_if=StoreType postgres"`
	// Sqlite configuration
	SqliteOpts SqliteConfig `json:"SQLITE" validate:"required_if=StoreType sqlite"`
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
			URI: os.Getenv("POSTGRES_URI"),
		},
		Encryption: EncryptionConfig{
			EncIV:     secretVals["ENCRYPTION_IV"],
			EncSecret: secretVals["ENCRYPTION_SECRET"],
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
