package services

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"math/big"

	"github.com/oalexander6/passman/pkg/entities"
)

func (s *Services) GetNoteByID(ctx context.Context, id entities.ID) (entities.Note, error) {
	note, err := s.noteStore.GetByID(ctx, id)
	if err != nil {
		return entities.Note{}, err
	}

	if note.Secure {
		decryptedVal, err := s.Decyrpt([]byte(note.Value))
		if err != nil {
			return entities.Note{}, entities.ErrDecryptFailed
		}

		note.Value = decryptedVal
	}

	return note, nil
}

func (s *Services) GetAllNotes(ctx context.Context) ([]entities.Note, error) {
	notes, err := s.noteStore.GetAll(ctx)
	if err != nil {
		return []entities.Note{}, err
	}

	for i, note := range notes {
		if note.Secure {
			decryptedVal, err := s.Decyrpt([]byte(note.Value))
			if err != nil {
				return []entities.Note{}, entities.ErrDecryptFailed
			}

			notes[i].Value = decryptedVal
		}
	}

	return notes, nil
}

func (s *Services) CreateNote(ctx context.Context, noteInput entities.NoteInput) (entities.Note, error) {
	if noteInput.Secure {
		encVal, err := s.Encrypt([]byte(noteInput.Value))
		if err != nil {
			return entities.Note{}, err
		}

		noteInput.Value = encVal
	}
	return s.noteStore.Create(ctx, noteInput)
}

func (s *Services) UpdateNote(ctx context.Context, note entities.Note) (entities.Note, error) {
	if note.Secure {
		encVal, err := s.Encrypt([]byte(note.Value))
		if err != nil {
			return entities.Note{}, err
		}

		note.Value = encVal
	}

	return s.noteStore.Update(ctx, note)
}

func (s *Services) DeleteNoteByID(ctx context.Context, id entities.ID) error {
	return s.noteStore.DeleteByID(ctx, id)
}

const validCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!.?#"

func (s *Services) GenerateRandomString(length int) (string, error) {
	result := make([]byte, length)

	for i := range result {
		charIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(validCharacters))))
		if err != nil {
			return "", err
		}

		result[i] = validCharacters[charIndex.Int64()]
	}

	return string(result), nil
}

// Encrypt implements AES-256 encryption using PKCS7 padding.
func (s *Services) Encrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher([]byte(s.config.EncKey))
	if err != nil {
		return "", err
	}

	paddedPlaintext := pkcs7Pad(plaintext, block.BlockSize())

	ciphertext := make([]byte, len(paddedPlaintext))

	mode := cipher.NewCBCEncrypter(block, []byte(s.config.EncIV))
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *Services) Decyrpt(encrypted []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(string(encrypted))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(s.config.EncKey))
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, []byte(s.config.EncIV))
	mode.CryptBlocks(ciphertext, ciphertext)

	plaintext, err := pkcs7UnPad(ciphertext, block.BlockSize())
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func pkcs7Pad(original []byte, padToMultipleOf int) []byte {
	ogLength := len(original)

	bytesToAdd := padToMultipleOf - ogLength%padToMultipleOf
	if bytesToAdd == 0 {
		bytesToAdd = padToMultipleOf
	}

	newBuf := make([]byte, ogLength+bytesToAdd)

	copy(newBuf, original)
	copy(newBuf[ogLength:], bytes.Repeat([]byte{uint8(bytesToAdd)}, bytesToAdd))

	return newBuf
}

func pkcs7UnPad(original []byte, blockSize int) ([]byte, error) {
	ogLength := len(original)
	if ogLength%blockSize != 0 {
		return []byte{}, errors.New("invalid padding")
	}

	bytesToRemove := int(original[ogLength-1])

	return original[:(ogLength - bytesToRemove)], nil
}
