package memory_store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/oalexander6/passman/pkg/entities"
	memory_store "github.com/oalexander6/passman/pkg/stores/memory"
)

func TestCreateNote(t *testing.T) {
	s := memory_store.New()

	noteInput := entities.NoteInput{
		Name:   "My Note",
		Value:  "Hello world",
		Secure: false,
	}

	result, err := s.NotesStore.Create(context.Background(), noteInput)
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if len(result.ID) == 0 {
		t.Fatal("Failed to generate an ID for the new note")
	}

	if result.Name != noteInput.Name || result.Value != noteInput.Value || result.Secure != noteInput.Secure {
		t.Fatal("Response from create contained incorrect values")
	}

	savedNote := s.Data.Notes[0]

	if savedNote.Name != noteInput.Name || savedNote.Value != noteInput.Value || savedNote.Secure != noteInput.Secure {
		t.Fatal("Saved note contained incorrect values")
	}
}

func TestDeleteNote(t *testing.T) {
	s := memory_store.New()

	note := entities.Note{
		ID:     uuid.NewString(),
		Name:   "My Note",
		Value:  "Hello world",
		Secure: false,
	}

	s.Data.Notes = append(s.Data.Notes, note)

	if err := s.NotesStore.DeleteByID(context.Background(), note.ID); err != nil {
		t.Fatal("Unexpected error")
	}

	if len(s.Data.Notes) != 0 {
		t.Fatal("Failed to delete note")
	}

	if err := s.NotesStore.DeleteByID(context.Background(), note.ID); err != nil {
		t.Fatal("Unexpected error")
	}
}

func TestGetAllNotes(t *testing.T) {
	s := memory_store.New()

	result, err := s.NotesStore.GetAll(context.Background())
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

	s.Data.Notes = []entities.Note{note1, note2}

	result, err = s.NotesStore.GetAll(context.Background())
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if len(result) != 2 {
		t.Fatal("Expected two results")
	}

	if !areEqualNotes(note1, result[0]) || !areEqualNotes(note2, result[1]) {
		t.Fatal("Results contain incorrect values")
	}
}

func TestGetNoteByID(t *testing.T) {
	s := memory_store.New()

	_, err := s.NotesStore.GetByID(context.Background(), uuid.NewString())
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

	s.Data.Notes = []entities.Note{note1, note2}

	result, err := s.NotesStore.GetByID(context.Background(), note1.ID)
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if !areEqualNotes(note1, result) {
		t.Fatal("Result contained incorrect values")
	}

	_, err = s.NotesStore.GetByID(context.Background(), uuid.NewString())
	if !errors.Is(err, entities.ErrNotFound) {
		t.Fatal("Expected not found error")
	}
}

func TestUpdateNote(t *testing.T) {
	s := memory_store.New()

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

	s.Data.Notes = []entities.Note{note1, note2}

	note1.Name = "Updated Note 1 Name"
	note1.Value = "Updated Note 1 Value"

	result, err := s.NotesStore.Update(context.Background(), note1)
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if !areEqualNotes(note1, result) {
		t.Fatal("Result contained incorrect values")
	}

	if !areEqualNotes(note1, s.Data.Notes[0]) {
		t.Fatal("Saved note contained incorrect values")
	}

	note1.ID = uuid.NewString()

	_, err = s.NotesStore.Update(context.Background(), note1)
	if !errors.Is(err, entities.ErrNotFound) {
		t.Fatal("Expected not found error")
	}
}

func areEqualNotes(note1 entities.Note, note2 entities.Note) bool {
	return note1.Name == note2.Name && note1.Value == note2.Value && note1.Secure == note2.Secure && note1.ID == note2.ID
}
