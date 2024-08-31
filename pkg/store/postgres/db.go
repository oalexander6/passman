package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oalexander6/passman/config"
	"github.com/oalexander6/passman/pkg/logger"
	"github.com/oalexander6/passman/pkg/models"
)

type PostgresStore struct {
	dbpool *pgxpool.Pool
}

var schema = `
CREATE TABLE IF NOT EXISTS notes (
	id         BIGSERIAL PRIMARY KEY,
	name       TEXT NOT NULL,
	value      TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL,
	deleted    BOOLEAN NOT NULL
);
`

func New(opts config.PostgresConfig) *PostgresStore {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := pgxpool.New(ctx, opts.URI)
	if err != nil {
		logger.Log.Fatal().Msgf("Unable to create pgx connection pool: %s", err)
	}

	if err = conn.Ping(ctx); err != nil {
		logger.Log.Fatal().Msgf("Failed to ping postgres: %s", err)
	}

	conn.Exec(context.Background(), schema)

	return &PostgresStore{
		dbpool: conn,
	}
}

func (s PostgresStore) Close() {
	s.dbpool.Close()
}

// NoteCreate implements models.Store.
func (s PostgresStore) NoteCreate(ctx context.Context, noteInput models.NoteCreateParams) (models.Note, error) {
	query := `INSERT INTO notes (name, value, created_at, updated_at, deleted) VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	currTime := time.Now().UTC().Format(time.RFC3339)

	var insertedID int64
	if err := s.dbpool.QueryRow(ctx, query, noteInput.Name, noteInput.Value, currTime, currTime, true).Scan(&insertedID); err != nil {
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

	result, err := s.dbpool.Exec(ctx, query, id)
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

	row, err := s.dbpool.Query(ctx, query, id)
	if err != nil {
		return models.Note{}, err
	}

	pgx.RowTo[models.Note](row)
	note, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Note])
	if err != nil {
		return models.Note{}, models.ErrNotFound
	}

	return note, nil
}

// NoteGetAll implements models.Store.
func (s PostgresStore) NoteGetAll(ctx context.Context) ([]models.Note, error) {
	query := `SELECT * FROM notes WHERE deleted=false;`

	rows, err := s.dbpool.Query(ctx, query)
	if err != nil {
		return []models.Note{}, err
	}

	notes, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Note])
	if err != nil {
		return []models.Note{}, models.ErrNotFound
	}

	return notes, nil
}
