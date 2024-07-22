package memory_store

import (
	"context"

	"github.com/google/uuid"
	"github.com/oalexander6/passman/pkg/entities"
)

type MemoryAccountsStore struct {
	Data *MemoryStoreData
}

// Create implements entities.AccountStore.
func (m *MemoryAccountsStore) Create(ctx context.Context, accountInput entities.AccountInput) (entities.Account, error) {
	id := uuid.NewString()

	account := entities.Account{
		ID:       id,
		Email:    accountInput.Email,
		Password: accountInput.Password,
	}

	m.Data.Accounts = append(m.Data.Accounts, account)

	return account, nil
}

// Delete implements entities.AccountStore.
func (m *MemoryAccountsStore) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// GetByID implements entities.AccountStore.
func (m *MemoryAccountsStore) GetByID(ctx context.Context, id string) (entities.Account, error) {
	for _, account := range m.Data.Accounts {
		if account.ID == id {
			return account, nil
		}
	}

	return entities.Account{}, entities.ErrNotFound
}

// GetByEmail implements entities.AccountStore.
func (m *MemoryAccountsStore) GetByEmail(ctx context.Context, email string) (entities.Account, error) {
	for _, account := range m.Data.Accounts {
		if account.Email == email {
			return account, nil
		}
	}

	return entities.Account{}, entities.ErrNotFound
}
