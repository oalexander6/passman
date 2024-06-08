package memory_store

import (
	"context"

	"github.com/google/uuid"
	"github.com/oalexander6/passman/pkg/entities"
)

type MemoryNotesStore struct {
	Data *MemoryStoreData
}

// Create implements entities.NoteStore.
func (m *MemoryNotesStore) Create(ctx context.Context, noteInput entities.NoteInput) (entities.Note, error) {
	id := uuid.NewString()

	note := entities.Note{
		ID:     id,
		Name:   noteInput.Name,
		Value:  noteInput.Value,
		Secure: noteInput.Secure,
	}

	m.Data.Notes = append(m.Data.Notes, note)

	return note, nil
}

// DeleteByID implements entities.NoteStore.
func (m *MemoryNotesStore) DeleteByID(ctx context.Context, id entities.ID) error {
	newNotes := make([]entities.Note, 0)

	for _, note := range m.Data.Notes {
		if note.ID != id {
			newNotes = append(newNotes, note)
		}
	}

	m.Data.Notes = newNotes

	return nil
}

// GetAll implements entities.NoteStore.
func (m *MemoryNotesStore) GetAll(ctx context.Context) ([]entities.Note, error) {
	return m.Data.Notes, nil
}

// GetByID implements entities.NoteStore.
func (m *MemoryNotesStore) GetByID(ctx context.Context, id entities.ID) (entities.Note, error) {
	for _, note := range m.Data.Notes {
		if note.ID == id {
			return note, nil
		}
	}

	return entities.Note{}, entities.ErrNotFound
}

// Update implements entities.NoteStore.
func (m *MemoryNotesStore) Update(ctx context.Context, note entities.Note) (entities.Note, error) {
	for index, storedNote := range m.Data.Notes {
		if storedNote.ID == note.ID {
			m.Data.Notes[index] = note
			return note, nil
		}
	}

	return entities.Note{}, entities.ErrNotFound
}
