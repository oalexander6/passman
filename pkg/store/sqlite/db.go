package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/oalexander6/passman/config"
	"github.com/oalexander6/passman/pkg/logger"
	"github.com/oalexander6/passman/pkg/models"
)

var dropAllTables = `
DROP TABLE accounts;
`

var schema = `
CREATE TABLE IF NOT EXISTS accounts (
	id         INTEGER PRIMARY KEY,
	email      TEXT NOT NULL UNIQUE,
	password   TEXT,
	name       TEXT NOT NULL,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	deleted    BOOLEAN NOT NULL
);
`

type SqliteStore struct {
	db *sqlx.DB
}

func New(opts config.SqliteConfig) *SqliteStore {
	db, err := sqlx.Connect("sqlite3", opts.DBFile)
	if err != nil {
		logger.Log.Fatal().Msgf("failed to connect to SQLite db: %s", err)
	}

	if opts.DeleteOnStartup {
		db.Exec(dropAllTables)
	}

	db.MustExec(schema)

	return &SqliteStore{
		db: db,
	}
}

// AccountCreate implements models.Store.
func (s *SqliteStore) AccountCreate(ctx context.Context, account models.Account) (models.Account, error) {
	baseFields := getNewBaseFields()

	account.CreatedAt = baseFields.CreatedAt
	account.UpdatedAt = baseFields.UpdatedAt
	account.Deleted = baseFields.Deleted

	query := `
	INSERT INTO accounts (email, password, name, created_at, updated_at, deleted)
	VALUES (:email, :password, :name, :created_at, :updated_at, :deleted);
	`

	result, err := s.db.NamedExecContext(ctx, query, account)
	if err != nil {
		return models.Account{}, err
	}

	accountID, err := result.LastInsertId()
	if err != nil {
		return models.Account{}, err
	}

	account.ID = accountID

	return account, nil
}

// AccountDelete implements models.Store.
func (s *SqliteStore) AccountDelete(ctx context.Context, id int) error {
	panic("unimplemented")
}

// AccountGetByEmail implements models.Store.
func (s *SqliteStore) AccountGetByEmail(ctx context.Context, email string) (models.Account, error) {
	query := `
	SELECT * FROM accounts WHERE email=$1;
	`

	var account models.Account

	if err := s.db.GetContext(ctx, &account, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Account{}, models.ErrNotFound
		}
		return models.Account{}, err
	}

	return account, nil
}

// AccountGetByID implements models.Store.
func (s *SqliteStore) AccountGetByID(ctx context.Context, id int) (models.Account, error) {
	panic("unimplemented")
}

// NoteCreate implements models.Store.
func (s *SqliteStore) NoteCreate(ctx context.Context, noteInput models.Note) (models.Note, error) {
	panic("unimplemented")
}

// NoteDeleteByID implements models.Store.
func (s *SqliteStore) NoteDeleteByID(ctx context.Context, id int) error {
	panic("unimplemented")
}

// NoteGetByAccountID implements models.Store.
func (s *SqliteStore) NoteGetByAccountID(ctx context.Context, userID int) ([]models.Note, error) {
	panic("unimplemented")
}

// NoteGetByID implements models.Store.
func (s *SqliteStore) NoteGetByID(ctx context.Context, id int) (models.Note, error) {
	panic("unimplemented")
}

// NoteUpdate implements models.Store.
func (s *SqliteStore) NoteUpdate(ctx context.Context, note models.Note) (models.Note, error) {
	panic("unimplemented")
}

func getNewBaseFields() models.Base {
	currTime := time.Now().UTC().Format(time.RFC3339)

	return models.Base{
		CreatedAt: currTime,
		UpdatedAt: currTime,
		Deleted:   false,
	}
}
