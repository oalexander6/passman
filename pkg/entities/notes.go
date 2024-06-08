package entities

import "context"

type Note struct {
	ID     ID
	Name   string
	Value  string
	Secure bool
}

type NoteInput struct {
	Name   string
	Value  string
	Secure bool
}

type NoteStore interface {
	GetByID(ctx context.Context, id ID) (Note, error)
	GetAll(ctx context.Context) ([]Note, error)
	Create(ctx context.Context, noteInput NoteInput) (Note, error)
	Update(ctx context.Context, note Note) (Note, error)
	DeleteByID(ctx context.Context, id ID) error
}
