package store

import (
	"context"

	"github.com/oalexander6/passman/pkg/models"
)

type Store struct {
}

// AccountDelete implements models.Store.
func (s Store) AccountDelete(ctx context.Context, id int) error {
	panic("unimplemented")
}

// AccountGetByEmail implements models.Store.
func (s Store) AccountGetByEmail(ctx context.Context, email string) (models.Account, error) {
	panic("unimplemented")
}

// NoteCreate implements models.Store.
func (s Store) NoteCreate(ctx context.Context, noteInput models.Note) (models.Note, error) {
	panic("unimplemented")
}

// NoteDeleteByID implements models.Store.
func (s Store) NoteDeleteByID(ctx context.Context, id int) error {
	panic("unimplemented")
}

// NoteGetByAccountID implements models.Store.
func (s Store) NoteGetByAccountID(ctx context.Context, userID int) ([]models.Note, error) {
	panic("unimplemented")
}

// NoteGetByID implements models.Store.
func (s Store) NoteGetByID(ctx context.Context, id int) (models.Note, error) {
	panic("unimplemented")
}

// NoteUpdate implements models.Store.
func (s Store) NoteUpdate(ctx context.Context, note models.Note) (models.Note, error) {
	panic("unimplemented")
}

// AccountCreate implements models.Store.
func (s Store) AccountCreate(ctx context.Context, account models.Account) (models.Account, error) {
	panic("unimplemented")
}

// AccountGetByID implements models.Store.
func (s Store) AccountGetByID(ctx context.Context, id int) (models.Account, error) {
	panic("unimplemented")
}
