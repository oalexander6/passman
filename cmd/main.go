package main

import (
	"os"

	"github.com/oalexander6/passman/config"
	"github.com/oalexander6/passman/httpserver"
	"github.com/oalexander6/passman/logger"
	"github.com/oalexander6/passman/models"
	"github.com/oalexander6/passman/store/postgres"
	"github.com/rs/zerolog"
)

func main() {
	logger.Init(zerolog.DebugLevel, os.Stdout)

	c := config.New()
	if err := c.Validate(); err != nil {
		logger.Log.Fatal().Msgf("Invalid configuration: %s", err.Error())
	}

	logger.Log.Info().Interface("config", c).Msg("Config initialized")

	var s models.Store

	switch c.StoreType {
	case config.STORE_TYPE_POSTGRES:
		s = postgres.New(c.PostgresOpts)
	default:
		logger.Log.Fatal().Msgf("Invalid store type: %s", c.StoreType)
	}

	defer s.Close()

	app := httpserver.New(c, *models.New(s, c.Encryption))

	app.Run()
}
