package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oalexander6/passman/config"
	"github.com/oalexander6/passman/pkg/logger"
	"github.com/oalexander6/passman/pkg/models"
)

type PostgresStore struct {
	dbpool *pgxpool.Pool
}

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

	return &PostgresStore{
		dbpool: conn,
	}
}

func (s PostgresStore) Close() {
	s.dbpool.Close()
}

// AccountDelete implements models.Store.
func (s PostgresStore) AccountDelete(ctx context.Context, id int64) error {
	panic("unimplemented")
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
	panic("unimplemented")
}

// NoteGetByAccountID implements models.Store.
func (s PostgresStore) NoteGetByAccountID(ctx context.Context, userID int64) ([]models.Note, error) {
	panic("unimplemented")
}

// NoteGetByID implements models.Store.
func (s PostgresStore) NoteGetByID(ctx context.Context, id int64) (models.Note, error) {
	panic("unimplemented")
}

// NoteUpdate implements models.Store.
func (s PostgresStore) NoteUpdate(ctx context.Context, note models.Note) (models.Note, error) {
	panic("unimplemented")
}

// AccountCreate implements models.Store.
func (s PostgresStore) AccountCreate(ctx context.Context, account models.Account) (models.Account, error) {
	panic("unimplemented")
}

// AccountGetByID implements models.Store.
func (s PostgresStore) AccountGetByID(ctx context.Context, id int64) (models.Account, error) {
	panic("unimplemented")
}
