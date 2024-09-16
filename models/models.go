package models

import "github.com/oalexander6/passman/config"

type Store interface {
	noteStore
	Close()
}

type Models struct {
	store          Store
	encyrptionOpts config.EncryptionConfig
}

func New(store Store, encryptionOpts config.EncryptionConfig) *Models {
	return &Models{
		store:          store,
		encyrptionOpts: encryptionOpts,
	}
}
