package models_test

import (
	"os"
	"testing"

	"github.com/oalexander6/passman/config"
	"github.com/oalexander6/passman/pkg/models"
	"github.com/oalexander6/passman/pkg/store/sqlite"
)

var m *models.Models

func TestMain(main *testing.M) {
	os.Remove("test.db")

	defer func() {
		os.Remove("test.db")
	}()

	c := &config.Config{
		Env:                  config.LOCAL_ENV,
		Port:                 "8080",
		Version:              "",
		SecretKey:            "test-secret",
		CSRFKey:              "test-csrf-key",
		EnableCSRFProtection: false,
		Encryption: config.EncryptionConfig{
			EncIV:     "not-the-real-iv!",
			EncSecret: "not-the-real-encryption-secret!!",
		},
		StoreType: config.STORE_TYPE_SQLITE,
		SqliteOpts: config.SqliteConfig{
			DBFile:          "test.db",
			DeleteOnStartup: true,
		},
	}

	m = models.New(sqlite.New(c.SqliteOpts), c)

	main.Run()
}
