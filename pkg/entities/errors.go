package entities

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("entity not found")
	ErrEncryptFailed = errors.New("encryption failed")
	ErrDecryptFailed = errors.New("decryption failed")
)
