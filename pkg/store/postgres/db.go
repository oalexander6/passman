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

var accountsSchema = `
CREATE TABLE IF NOT EXISTS accounts (
	id         BIGSERIAL PRIMARY KEY,
	email      TEXT NOT NULL,
	password   TEXT NOT NULL,
	name       TEXT NOT NULL,
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

	var greeting string
	err = conn.QueryRow(context.Background(), "SELECT 'Hello, world!'").Scan(&greeting)
	if err != nil {
		logger.Log.Fatal().Msgf("Failed to get greeting: %s", err)
	}

	logger.Log.Debug().Msgf("Got greeting: %s", greeting)

	conn.Exec(context.Background(), accountsSchema)

	return &PostgresStore{
		dbpool: conn,
	}
}

func (s PostgresStore) Close() {
	s.dbpool.Close()
}

// AccountDelete implements models.Store.
func (s PostgresStore) AccountDelete(ctx context.Context, id int64) error {
	query := `UPDATE accounts SET deleted=true WHERE id=$1;`

	result, err := s.dbpool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() != 1 {
		return models.ErrNotFound
	}

	return nil
}

// AccountGetByEmail implements models.Store.
func (s PostgresStore) AccountGetByEmail(ctx context.Context, email string) (models.Account, error) {
	panic("unimplemented")
}

// NoteCreate implements models.Store.
func (s PostgresStore) NoteCreate(ctx context.Context, noteInput models.Note) (models.Note, error) {
	panic("unimplemented")
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

// NoteGetByAccountID implements models.Store.
func (s PostgresStore) NoteGetByAccountID(ctx context.Context, userID int64) ([]models.Note, error) {
	panic("unimplemented")
}

// NoteGetByID implements models.Store.
func (s PostgresStore) NoteGetByID(ctx context.Context, id int64) (models.Note, error) {
	query := `SELECT * FROM notes WHERE id=$1 AND deleted=false;`

	row, err := s.dbpool.Query(ctx, query, id)
	if err != nil {
		return models.Note{}, err
	}

	pgx.RowTo[models.Account](row)
	account, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Account])
	if err != nil {
		return models.Account{}, models.ErrNotFound
	}

	return account, nil
}

// NoteUpdate implements models.Store.
func (s PostgresStore) NoteUpdate(ctx context.Context, note models.Note) (models.Note, error) {
	panic("unimplemented")
}

// AccountCreate implements models.Store.
func (s PostgresStore) AccountCreate(ctx context.Context, account models.Account) (models.Account, error) {
	query := `INSERT INTO accounts (email, password, name, created_at, updated_at, delete) VALUES (@email, @password, @name, @created_at, @updated_at, @delete);`

	args := pgx.NamedArgs{
		"email":      account.Email,
		"password":   account.Password,
		"name":       account.Name,
		"created_at": time.Now().UTC().Format(time.RFC3339),
		"updated_at": time.Now().UTC().Format(time.RFC3339),
		"deleted":    false,
	}

	s.dbpool.Exec(ctx, query, args)

	return models.Account{}, nil
}

// AccountGetByID implements models.Store.
func (s PostgresStore) AccountGetByID(ctx context.Context, id int64) (models.Account, error) {
	query := `SELECT * FROM accounts WHERE id=$1 AND deleted=false;`

	row, err := s.dbpool.Query(ctx, query, id)
	if err != nil {
		return models.Account{}, err
	}

	pgx.RowTo[models.Account](row)
	account, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Account])
	if err != nil {
		return models.Account{}, models.ErrNotFound
	}

	return account, nil
}
