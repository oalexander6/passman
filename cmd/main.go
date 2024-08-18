package main

import (
	"os"

	"github.com/oalexander6/passman/config"
	"github.com/oalexander6/passman/pkg/logger"
	"github.com/oalexander6/passman/pkg/server"
	"github.com/oalexander6/passman/pkg/store"
	"github.com/rs/zerolog"
)

func main() {
	logger.Init(zerolog.DebugLevel, os.Stdout)

	c := config.New()
	if err := c.Validate(); err != nil {
		logger.Log.Fatal().Msgf("Invalid configuration: %s", err.Error())
	}

	app := server.New(c, store.Store{})
	logger.Log.Fatal().Msgf("Application crashed: %s", app.Run().Error())
}
