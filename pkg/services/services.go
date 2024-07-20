package services

import (
	"github.com/oalexander6/passman/pkg/config"
	"github.com/oalexander6/passman/pkg/entities"
)

// Services represents the services available to the application that implement
// functionality.
type Services struct {
	noteStore     entities.NoteStore
	accountsStore entities.AccountStore
	config        *config.Config
}

func New(config *config.Config, noteStore entities.NoteStore, accountStore entities.AccountStore) *Services {
	return &Services{
		noteStore:     noteStore,
		accountsStore: accountStore,
		config:        config,
	}
}
