package models

import "errors"

var (
	ErrNotFound      = errors.New("entity not found")
	ErrAlreadyExists = errors.New("entity already exists")
	ErrEncryptFailed = errors.New("encryption failed")
	ErrDecryptFailed = errors.New("decryption failed")
)
