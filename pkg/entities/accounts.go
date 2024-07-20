package entities

import "context"

type Account struct {
	ID ID `json:"id"`
}

type AccountStore interface {
	GetByID(ctx context.Context, id ID) (Account, error)
	Create(ctx context.Context) (Account, error)
	Delete(ctx context.Context, id ID) error
}
