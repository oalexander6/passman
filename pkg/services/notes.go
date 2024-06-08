package services

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"math/big"

	"github.com/oalexander6/passman/pkg/entities"
)

// GetNoteByID returns the note with the provided ID with the value of secure notes decrypted.
// Returns an error if the note is not found.
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

// GetAllNotes returns all stored notes with the value of secure notes decrypted.
// Does NOT return an error if no notes are found.
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

// CreateNote saves a new note. It will encrypt the value of the note if it is marked as secure.
// Returns an error if the note fails to save.
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

// UpdateNote updates the note that matches the provided note's ID. It will encrypt the provided
// value if the note is marked as secure.
// Returns an error if no note with the provided ID is found.
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

// DeleteNoteByID will remove the note with the provided ID.
// Returns an error if a note with that ID is not found.
func (s *Services) DeleteNoteByID(ctx context.Context, id entities.ID) error {
	return s.noteStore.DeleteByID(ctx, id)
}

// GenerateRandomString returns a cryptographically secure random string of the provided length.
func (s *Services) GenerateRandomString(length int) (string, error) {
	result := make([]byte, length)

	const validCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!.?#"

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

// Decrypt implements AES-256 decryption using PKCS7 unpadding.
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

// pkcs7Pad implements PKCS7 padding by checking the length of the provided
// buffer, and adding the number of bytes required to increase the length to
// the next multiple of padToMultipleOf. It always pads with at least one byte.
// The inserted bytes are all set to the number of bytes inserted.
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

// pkcs7UnPad implements removal of PKCS7 padding by checking the value of the
// last byte and removing that many bytes from the end of the original buffer.
func pkcs7UnPad(original []byte, blockSize int) ([]byte, error) {
	ogLength := len(original)
	if ogLength%blockSize != 0 {
		return []byte{}, entities.ErrDecryptFailed
	}

	bytesToRemove := int(original[ogLength-1])

	return original[:(ogLength - bytesToRemove)], nil
}
