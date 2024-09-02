package models

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"math/big"
)

// Note represents a note/password. The value field will always be stored encrypted.
type Note struct {
	ID        int64
	Name      string
	Value     string
	CreatedAt string
	UpdatedAt string
	Deleted   bool
}

// NoteCreateParams represents the data required to create a new note.
type NoteCreateParams struct {
	Name  string `json:"name" form:"name" binding:"required"`
	Value string `json:"value" form:"value" binding:"required"`
}

// NoteCreateRandomParams represents the data required to create a new random note.
type NoteCreateRandomParams struct {
	Name   string `json:"name" form:"name" binding:"required"`
	Length int    `json:"length" form:"length" binding:"required,lte=2048"`
}

// NoteGetResponse represents the data returned for note GET requests.
type NoteGetResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// NoteStore defines the interface required to implement persistent storage functionality
// for notes.
type noteStore interface {
	NoteGetByID(ctx context.Context, id int64) (Note, error)
	NoteGetAll(ctx context.Context) ([]Note, error)
	NoteCreate(ctx context.Context, noteInput NoteCreateParams) (Note, error)
	NoteDeleteByID(ctx context.Context, id int64) error
}

// NoteGetByID returns the note with the provided ID with the value decrypted.
// Returns an error if the note is not found.
func (m *Models) NoteGetByID(ctx context.Context, noteID int64) (NoteGetResponse, error) {
	note, err := m.store.NoteGetByID(ctx, noteID)
	if err != nil {
		return NoteGetResponse{}, err
	}

	decryptedVal, err := m.Decyrpt([]byte(note.Value))
	if err != nil {
		return NoteGetResponse{}, ErrDecryptFailed
	}

	note.Value = decryptedVal

	return NoteGetResponse{
		Name:  note.Name,
		Value: note.Value,
	}, nil
}

// NoteGetAll returns all notes with their value's decrypted.
// Returns an error if the note is not found.
func (m *Models) NoteGetAll(ctx context.Context) ([]NoteGetResponse, error) {
	notes, err := m.store.NoteGetAll(ctx)
	if err != nil {
		return []NoteGetResponse{}, err
	}

	results := make([]NoteGetResponse, len(notes))
	for i := range notes {
		decryptedVal, err := m.Decyrpt([]byte(notes[i].Value))
		if err != nil {
			return []NoteGetResponse{}, ErrDecryptFailed
		}

		results[i] = NoteGetResponse{
			ID:        notes[i].ID,
			Name:      notes[i].Name,
			Value:     decryptedVal,
			CreatedAt: notes[i].CreatedAt,
			UpdatedAt: notes[i].UpdatedAt,
		}
	}

	return results, nil
}

// NoteCreate saves a new note. It will encrypt the value of the note if it is marked as secure.
// Returns an error if the note fails to save.
func (m *Models) NoteCreate(ctx context.Context, noteInput NoteCreateParams) (NoteGetResponse, error) {
	encVal, err := m.Encrypt([]byte(noteInput.Value))
	if err != nil {
		return NoteGetResponse{}, err
	}

	noteInput.Value = encVal

	savedNote, err := m.store.NoteCreate(ctx, noteInput)
	if err != nil {
		return NoteGetResponse{}, err
	}

	decryptedVal, err := m.Decyrpt([]byte(savedNote.Value))
	if err != nil {
		return NoteGetResponse{}, ErrDecryptFailed
	}

	return NoteGetResponse{
		ID:        savedNote.ID,
		Name:      savedNote.Name,
		Value:     decryptedVal,
		CreatedAt: savedNote.CreatedAt,
		UpdatedAt: savedNote.UpdatedAt,
	}, nil
}

// NoteCreateRandom saves a new note with a randomly generated value.
// Returns an error if the note fails to save or the random value generation fails.
func (m *Models) NoteCreateRandom(ctx context.Context, noteInput NoteCreateRandomParams) (NoteGetResponse, error) {
	validCharacters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!?.#$"
	randomVal, err := generateRandomString(noteInput.Length, validCharacters)
	if err != nil {
		return NoteGetResponse{}, err
	}

	noteCreateParams := NoteCreateParams{
		Name:  noteInput.Name,
		Value: randomVal,
	}

	return m.NoteCreate(ctx, noteCreateParams)
}

// DeleteNoteByID will remove the note with the provided ID.
// Returns an error if a note with that ID is not found.
func (m *Models) NoteDeleteByID(ctx context.Context, noteID int64) error {
	return m.store.NoteDeleteByID(ctx, noteID)
}

// generateRandomString returns a cryptographically secure random string of the provided length.
func generateRandomString(length int, validCharacters string) (string, error) {
	if len(validCharacters) == 0 {
		return "", errors.New("must provide at least one valid character")
	}

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
func (m *Models) Encrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher([]byte(m.config.Encryption.EncSecret))
	if err != nil {
		return "", err
	}

	paddedPlaintext := pkcs7Pad(plaintext, block.BlockSize())

	ciphertext := make([]byte, len(paddedPlaintext))

	mode := cipher.NewCBCEncrypter(block, []byte(m.config.Encryption.EncIV))
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt implements AES-256 decryption using PKCS7 unpadding.
func (m *Models) Decyrpt(encrypted []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(string(encrypted))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(m.config.Encryption.EncSecret))
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, []byte(m.config.Encryption.EncIV))
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
		return []byte{}, ErrDecryptFailed
	}

	bytesToRemove := int(original[ogLength-1])

	return original[:(ogLength - bytesToRemove)], nil
}
