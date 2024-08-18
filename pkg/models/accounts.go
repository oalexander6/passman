package models

import (
	"context"
	"errors"

	"github.com/alexedwards/argon2id"
)

// Account represents a user account of any type. An account may be stored with an
// empty string for a password, indicating that they must log in with OAuth.
type Account struct {
	ID       int    `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Name     string `db:"name"`
	base
}

// Represents the required input to the register method.
type AccountCreateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// Represents the required input to the login method.
type AccountLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Represents the type of the response from the get all and get one methods.
type AccountGetResponse struct {
	ID    int    `db:"id"`
	Email string `db:"email"`
	Name  string `db:"name"`
}

// Defines the required interface to implement an account store.
type accountStore interface {
	AccountCreate(ctx context.Context, account Account) (Account, error)
	AccountGetByID(ctx context.Context, id int) (Account, error)
	AccountGetByEmail(ctx context.Context, email string) (Account, error)
	AccountDelete(ctx context.Context, id int) error
}

// Checks for existing account, creates a new account and saves it with the password hashed.
func (m *Models) AccountRegister(ctx context.Context, account AccountCreateRequest) (IDResponse, error) {
	_, err := m.store.AccountGetByEmail(ctx, account.Email)
	if err == nil {
		return IDResponse{}, ErrAlreadyExists
	}
	if !errors.Is(err, ErrNotFound) {
		return IDResponse{}, err
	}

	hashedPassword, err := argon2id.CreateHash(account.Password, argon2id.DefaultParams)
	if err != nil {
		return IDResponse{}, errors.New("password hash failed")
	}

	accountToStore := Account{
		Name:     account.Name,
		Email:    account.Email,
		Password: hashedPassword,
	}

	savedAccount, err := m.store.AccountCreate(ctx, accountToStore)
	if err != nil {
		return IDResponse{}, err
	}

	return IDResponse{
		ID: savedAccount.ID,
	}, nil
}

// Checks the provided credentials against the existing account and it's password hash. Returns an error
// if the credentials are invalid.
func (m *Models) AccountLogin(ctx context.Context, credentials AccountLoginRequest) (IDResponse, error) {
	return IDResponse{}, nil
}

// Retrieves the account with the provided ID. Returns an error if the ID is not found.
func (m *Models) AccountGetByID(ctx context.Context, id int) (AccountGetResponse, error) {
	return AccountGetResponse{}, nil
}

// Retrieves all accounts. Returns an error if none are found.
func (m *Models) AccountGetAll(ctx context.Context) ([]AccountGetResponse, error) {
	return []AccountGetResponse{}, nil
}
