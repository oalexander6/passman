package entities

import "context"

type Account struct {
	ID       ID     `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type AccountInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AccountLoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AccountStore interface {
	GetByID(ctx context.Context, id ID) (Account, error)
	GetByEmail(ctx context.Context, email string) (Account, error)
	Create(ctx context.Context, accountInput AccountInput) (Account, error)
	Delete(ctx context.Context, id ID) error
}
