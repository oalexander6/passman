package entities

import "context"

// Note represents a note/password, which may be secure or not secure. Secure notes will
// have their value encrypted upon storage.
type Note struct {
	ID     ID     `json:"id" form:"id" binding="required"`
	Name   string `json:"name" form:"name" binding="required"`
	Value  string `json:"value" form:"value" binding="required"`
	Secure bool   `json:"secure" form:"secure"`
}

// NoteInput represents the data required to create a new note.
type NoteInput struct {
	Name   string `json:"name" form:"name" binding:"required"`
	Value  string `json:"value" form:"value" binding:"required"`
	Secure bool   `json:"secure" form:"secure"`
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
