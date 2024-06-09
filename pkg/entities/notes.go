package entities

import "context"

// Note represents a note/password, which may be secure or not secure. Secure notes will
// have their value encrypted upon storage.
type Note struct {
	ID     ID     `json:"id"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	Secure bool   `json:"secure"`
}

// NoteInput represents the data required to create a new note.
type NoteInput struct {
	Name   string `json:"name" binding:"required"`
	Value  string `json:"value" binding:"required"`
	Secure bool   `json:"secure"`
}

// NoteStore defines the interface required to implement persistent storage functionality
// for notes.
type NoteStore interface {
	GetByID(ctx context.Context, id ID) (Note, error)
	GetAll(ctx context.Context) ([]Note, error)
	Create(ctx context.Context, noteInput NoteInput) (Note, error)
	Update(ctx context.Context, note Note) (Note, error)
	DeleteByID(ctx context.Context, id ID) error
}
