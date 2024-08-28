package main

import (
	"os"

	"github.com/oalexander6/passman/config"
	"github.com/oalexander6/passman/pkg/httpserver"
	"github.com/oalexander6/passman/pkg/logger"
	"github.com/oalexander6/passman/pkg/models"
	"github.com/oalexander6/passman/pkg/store/postgres"
	"github.com/rs/zerolog"
)

func main() {
	logger.Init(zerolog.DebugLevel, os.Stdout)

	c := config.New()
	if err := c.Validate(); err != nil {
		logger.Log.Fatal().Msgf("Invalid configuration: %s", err.Error())
	}

	var store models.Store

	switch c.StoreType {
	case config.STORE_TYPE_POSTGRES:
		store = postgres.New(c.PostgresOpts)
	default:
		logger.Log.Fatal().Msgf("Invalid store type: %s", c.StoreType)
	}

	defer store.Close()

	app := httpserver.New(c, store)
	logger.Log.Fatal().Msgf("Application crashed: %s", app.Run().Error())
}
