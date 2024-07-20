package services_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/oalexander6/passman/pkg/config"
	"github.com/oalexander6/passman/pkg/entities"
	"github.com/oalexander6/passman/pkg/services"
	memory_store "github.com/oalexander6/passman/pkg/stores/memory"
)

var testConfig = &config.Config{
	EncIV:  "1234567890123456",
	EncKey: "12345678901234567890123456789012",
}

func TestGenerateRandomString(t *testing.T) {
	s := services.New(testConfig, nil, nil)

	const validCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!.?#"

	val, err := s.GenerateRandomString(12, validCharacters)
	if err != nil {
		t.FailNow()
	}

	if len(val) != 12 {
		t.Fatal("Random string of length 12 was expected")
	}

	newVal, err := s.GenerateRandomString(12, validCharacters)
	if err != nil {
		t.FailNow()
	}

	if newVal == val {
		t.Fatal("Expected different random string from subsequent call")
	}

	for _, char := range newVal {
		if !strings.ContainsRune(validCharacters, char) {
			t.Fatalf("Invalid character found: %c", char)
		}
	}

	longVal, err := s.GenerateRandomString(256, validCharacters)
	if err != nil {
		t.FailNow()
	}

	if len(longVal) != 256 {
		t.Fatal("Random string of length 256 was expected")
	}

	for _, char := range newVal {
		if !strings.ContainsRune(validCharacters, char) {
			t.Fatalf("Invalid character found: %c", char)
		}
	}

	if _, err = s.GenerateRandomString(10, ""); err == nil {
		t.Fatal("Expected 'must provide at least one valid character' error")
	}
}

func TestGetNoteByID(t *testing.T) {
	m := memory_store.New()
	s := services.New(testConfig, m.NotesStore, m.AccountsStore)

	_, err := s.GetNoteByID(context.Background(), uuid.NewString())
	if !errors.Is(err, entities.ErrNotFound) {
		t.Fatal("Expected not found error")
	}

	note1 := entities.Note{
		ID:     uuid.NewString(),
		Name:   "My Note 1",
		Value:  "Hello world 1",
		Secure: false,
	}

	note2 := entities.Note{
		ID:     uuid.NewString(),
		Name:   "My Note 2",
		Value:  "Hello world 2",
		Secure: false,
	}

	m.Data.Notes = []entities.Note{note1, note2}

	result, err := s.GetNoteByID(context.Background(), note1.ID)
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if !areEqualNotes(note1, result) {
		t.Fatal("Result contained incorrect values")
	}

	_, err = s.GetNoteByID(context.Background(), uuid.NewString())
	if !errors.Is(err, entities.ErrNotFound) {
		t.Fatal("Expected not found error")
	}

	unencryptedVal := "Hello world 3"
	encryptedVal, err := s.Encrypt([]byte(unencryptedVal))
	if err != nil {
		t.Fatal("Unexpected error")
	}

	secureNote := entities.Note{
		ID:     uuid.NewString(),
		Name:   "My Note",
		Value:  encryptedVal,
		Secure: true,
	}

	m.Data.Notes = append(m.Data.Notes, secureNote)

	result, err = s.GetNoteByID(context.Background(), secureNote.ID)
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if result.ID != secureNote.ID || result.Name != secureNote.Name ||
		result.Secure != secureNote.Secure || result.Value != unencryptedVal {
		t.Fatal("Failed to return correct secure note with decrypted value")
	}
}

func TestGetAllNotes(t *testing.T) {
	m := memory_store.New()
	s := services.New(testConfig, m.NotesStore, m.AccountsStore)

	result, err := s.GetAllNotes(context.Background())
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if len(result) != 0 {
		t.Fatal("Expected empty results")
	}

	note1 := entities.Note{
		ID:     uuid.NewString(),
		Name:   "My Note 1",
		Value:  "Hello world 1",
		Secure: false,
	}

	note2 := entities.Note{
		ID:     uuid.NewString(),
		Name:   "My Note 2",
		Value:  "Hello world 2",
		Secure: false,
	}

	m.Data.Notes = []entities.Note{note1, note2}

	result, err = s.GetAllNotes(context.Background())
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if len(result) != 2 {
		t.Fatal("Expected two results")
	}

	if !areEqualNotes(note1, result[0]) || !areEqualNotes(note2, result[1]) {
		t.Fatal("Results contain incorrect values")
	}

	unencryptedVal := "Hello world 3"
	encryptedVal, err := s.Encrypt([]byte(unencryptedVal))
	if err != nil {
		t.Fatal("Unexpected error")
	}

	secureNote := entities.Note{
		ID:     uuid.NewString(),
		Name:   "My Note",
		Value:  encryptedVal,
		Secure: true,
	}

	m.Data.Notes = append(m.Data.Notes, secureNote)

	result, err = s.GetAllNotes(context.Background())
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if len(result) != 3 {
		t.Fatal("Expected two results")
	}

	if result[2].ID != secureNote.ID || result[2].Name != secureNote.Name ||
		result[2].Secure != secureNote.Secure || result[2].Value != unencryptedVal {
		t.Fatal("Failed to return correct secure note with decrypted value")
	}
}

func TestCreateNote(t *testing.T) {
	m := memory_store.New()
	s := services.New(testConfig, m.NotesStore, m.AccountsStore)

	noteInput := entities.NoteInput{
		Name:   "My Note",
		Value:  "Hello world",
		Secure: false,
	}

	result, err := s.CreateNote(context.Background(), noteInput)
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if len(result.ID) == 0 {
		t.Fatal("Failed to generate an ID for the new note")
	}

	if result.Name != noteInput.Name || result.Value != noteInput.Value || result.Secure != noteInput.Secure {
		t.Fatal("Response from create contained incorrect values")
	}

	savedNote := m.Data.Notes[0]

	if savedNote.Name != noteInput.Name || savedNote.Value != noteInput.Value || savedNote.Secure != noteInput.Secure {
		t.Fatal("Saved note contained incorrect values")
	}

	secureNoteInput := entities.NoteInput{
		Name:   "My Note",
		Value:  "Hello world",
		Secure: true,
	}

	result, err = s.CreateNote(context.Background(), secureNoteInput)
	if err != nil {
		t.Fatal("Unexpected error")
	}

	encryptedVal, err := s.Encrypt([]byte(secureNoteInput.Value))
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if result.Value != secureNoteInput.Value {
		t.Fatal("Result expected to contain unencrypted value")
	}

	if m.Data.Notes[1].Value != encryptedVal {
		t.Fatal("Saved note expected to contain encrypted value")
	}
}

func TestDeleteNote(t *testing.T) {
	m := memory_store.New()
	s := services.New(testConfig, m.NotesStore, m.AccountsStore)

	note := entities.Note{
		ID:     uuid.NewString(),
		Name:   "My Note",
		Value:  "Hello world",
		Secure: false,
	}

	m.Data.Notes = append(m.Data.Notes, note)

	if err := s.DeleteNoteByID(context.Background(), note.ID); err != nil {
		t.Fatal("Unexpected error")
	}

	if len(m.Data.Notes) != 0 {
		t.Fatal("Failed to delete note")
	}

	if err := s.DeleteNoteByID(context.Background(), note.ID); err != nil {
		t.Fatal("Unexpected error")
	}
}

func TestUpdateNote(t *testing.T) {
	m := memory_store.New()
	s := services.New(testConfig, m.NotesStore, m.AccountsStore)

	note1 := entities.Note{
		ID:     uuid.NewString(),
		Name:   "My Note 1",
		Value:  "Hello world 1",
		Secure: false,
	}

	note2 := entities.Note{
		ID:     uuid.NewString(),
		Name:   "My Note 2",
		Value:  "Hello world 2",
		Secure: false,
	}

	m.Data.Notes = []entities.Note{note1, note2}

	note1.Name = "Updated Note 1 Name"
	note1.Value = "Updated Note 1 Value"

	result, err := s.UpdateNote(context.Background(), note1)
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if !areEqualNotes(note1, result) {
		t.Fatal("Result contained incorrect values")
	}

	if !areEqualNotes(note1, m.Data.Notes[0]) {
		t.Fatal("Saved note contained incorrect values")
	}

	note1.ID = uuid.NewString()

	_, err = s.UpdateNote(context.Background(), note1)
	if !errors.Is(err, entities.ErrNotFound) {
		t.Fatal("Expected not found error")
	}

	note2.Secure = true

	result, err = s.UpdateNote(context.Background(), note2)
	if err != nil {
		t.Fatal("Unexpected error")
	}

	encryptedVal, err := s.Encrypt([]byte(note2.Value))
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if result.Value != encryptedVal {
		t.Fatal("Result expected to contain encrypted value")
	}

	if m.Data.Notes[1].Value != encryptedVal {
		t.Fatal("Saved note expected to contain encrypted value")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	original := "This is my 32 byte string 123456"

	s := services.New(testConfig, nil, nil)

	ciphertext, err := s.Encrypt([]byte(original))
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if ciphertext == original {
		t.Fatal("Expected a different value than the original text")
	}

	decrypted, err := s.Decyrpt([]byte(ciphertext))
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if decrypted != original {
		t.Fatal("Decrypted value did not match original")
	}

	originalNonBlockSizeMultiple := "This is a string of a length that is not a multiple of the block size"

	ciphertext, err = s.Encrypt([]byte(originalNonBlockSizeMultiple))
	if err != nil {
		t.Fatal("Unexpected error")
	}

	decrypted, err = s.Decyrpt([]byte(ciphertext))
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if decrypted != originalNonBlockSizeMultiple {
		t.Fatal("Decrypted value did not match original")
	}
}

func areEqualNotes(note1 entities.Note, note2 entities.Note) bool {
	return note1.Name == note2.Name && note1.Value == note2.Value && note1.Secure == note2.Secure && note1.ID == note2.ID
}
