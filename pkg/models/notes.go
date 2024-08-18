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

// Note represents a note/password, which may be secure or not secure. Secure notes will
// have their value encrypted upon storage.
type Note struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Value string `db:"value"`
	base
}

// NoteCreateRequest represents the data required to create a new note.
type NoteCreateRequest struct {
	Name  string `json:"name" form:"name" validate:"required"`
	Value string `json:"value" form:"value" validate:"required"`
}

type NoteGetResponse struct {
	Name  string `json:"name" form:"name"`
	Value string `json:"value" form:"value"`
}

// NoteStore defines the interface required to implement persistent storage functionality
// for notes.
type noteStore interface {
	NoteGetByID(ctx context.Context, id int) (Note, error)
	NoteGetByAccountID(ctx context.Context, userID int) ([]Note, error)
	NoteCreate(ctx context.Context, noteInput Note) (Note, error)
	NoteUpdate(ctx context.Context, note Note) (Note, error)
	NoteDeleteByID(ctx context.Context, id int) error
}

// GetNoteByID returns the note with the provided ID with the value of secure notes decrypted.
// Returns an error if the note is not found.
func (m *Models) GetNoteByID(ctx context.Context, noteID int) (NoteGetResponse, error) {
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

// NoteGetByUserID returns all stored notes with the value of secure notes decrypted.
// Does NOT return an error if no notes are found.
func (m *Models) NoteGetByUserID(ctx context.Context, accountID int) ([]NoteGetResponse, error) {
	notes, err := m.store.NoteGetByAccountID(ctx, accountID)
	if err != nil {
		return []NoteGetResponse{}, err
	}

	unencryptedNotes := make([]NoteGetResponse, len(notes))

	for i := range notes {
		unencryptedNotes[i] = NoteGetResponse{
			Name: notes[i].Name,
		}

		decryptedVal, err := m.Decyrpt([]byte(notes[i].Value))
		if err != nil {
			return []NoteGetResponse{}, ErrDecryptFailed
		}

		unencryptedNotes[i].Value = decryptedVal
	}

	return unencryptedNotes, nil
}

// NoteCreate saves a new note. It will encrypt the value of the note if it is marked as secure.
// Returns an error if the note fails to save.
func (m *Models) NoteCreate(ctx context.Context, noteInput Note) (Note, error) {
	unencryptedVal := noteInput.Value

	encVal, err := m.Encrypt([]byte(noteInput.Value))
	if err != nil {
		return Note{}, err
	}

	noteInput.Value = encVal

	savedNote, err := m.store.NoteCreate(ctx, noteInput)
	if err != nil {
		return Note{}, err
	}

	savedNote.Value = unencryptedVal

	return savedNote, nil
}

// NoteUpdate updates the note that matches the provided note's ID. It will encrypt the provided
// value if the note is marked as secure.
// Returns an error if no note with the provided ID is found.
func (m *Models) NoteUpdate(ctx context.Context, note Note) (Note, error) {
	unencryptedVal := note.Value

	encVal, err := m.Encrypt([]byte(note.Value))
	if err != nil {
		return Note{}, err
	}

	note.Value = encVal
	savedNote, err := m.store.NoteUpdate(ctx, note)
	if err != nil {
		return Note{}, err
	}

	savedNote.Value = unencryptedVal

	return savedNote, nil
}

// DeleteNoteByID will remove the note with the provided ID.
// Returns an error if a note with that ID is not found.
func (m *Models) NoteDeleteByID(ctx context.Context, noteID int) error {
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
