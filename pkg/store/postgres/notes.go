package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oalexander6/passman/pkg/models"
)

type note struct {
	ID        int64              `db:"id"`
	Name      string             `db:"name"`
	Value     string             `db:"value"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
	Deleted   pgtype.Bool        `db:"deleted"`
}

// NoteCreate implements models.Store.
func (s PostgresStore) NoteCreate(ctx context.Context, noteInput models.NoteCreateParams) (models.Note, error) {
	query := `INSERT INTO notes (name, value, created_at, updated_at, deleted) VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	currTime := time.Now().UTC().Format(time.RFC3339)

	var insertedID int64
	if err := s.db.QueryRow(ctx, query, noteInput.Name, noteInput.Value, currTime, currTime, false).Scan(&insertedID); err != nil {
		return models.Note{}, err
	}

	return models.Note{
		ID:        insertedID,
		Name:      noteInput.Name,
		Value:     noteInput.Value,
		CreatedAt: currTime,
		UpdatedAt: currTime,
	}, nil
}

// NoteDeleteByID implements models.Store.
func (s PostgresStore) NoteDeleteByID(ctx context.Context, id int64) error {
	query := `UPDATE notes SET deleted=true WHERE id=$1;`

	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() != 1 {
		return models.ErrNotFound
	}

	return nil
}

// NoteGetByID implements models.Store.
func (s PostgresStore) NoteGetByID(ctx context.Context, id int64) (models.Note, error) {
	query := `SELECT * FROM notes WHERE id=$1 AND deleted=false;`

	row, err := s.db.Query(ctx, query, id)
	if err != nil {
		return models.Note{}, err
	}

	pgx.RowTo[note](row)
	note, err := pgx.CollectOneRow(row, pgx.RowToStructByName[note])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Note{}, models.ErrNotFound
		}
		return models.Note{}, err
	}

	return noteToModel(note), nil
}

// NoteGetAll implements models.Store.
func (s PostgresStore) NoteGetAll(ctx context.Context) ([]models.Note, error) {
	query := `SELECT * FROM notes WHERE deleted=false;`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return []models.Note{}, err
	}

	notes, err := pgx.CollectRows(rows, pgx.RowToStructByName[note])
	if err != nil {
		return []models.Note{}, err
	}

	return notesToModel(notes), nil
}

// Converts a DB note struct to a models.Note struct.
func noteToModel(note note) models.Note {
	return models.Note{
		ID:        note.ID,
		Name:      note.Name,
		Value:     note.Value,
		CreatedAt: note.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt: note.UpdatedAt.Time.Format(time.RFC3339),
		Deleted:   note.Deleted.Bool,
	}
}

// // Converts a list of DB note structs to a list of models.Note structs.
func notesToModel(notes []note) []models.Note {
	results := make([]models.Note, len(notes))

	for i := range notes {
		results[i] = noteToModel(notes[i])
	}

	return results
}
