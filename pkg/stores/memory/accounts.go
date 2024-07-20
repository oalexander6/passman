package memory_store

import (
	"context"

	"github.com/oalexander6/passman/pkg/entities"
)

type MemoryAccountsStore struct {
	Data *MemoryStoreData
}

// Create implements entities.AccountStore.
func (m *MemoryAccountsStore) Create(ctx context.Context) (entities.Account, error) {
	panic("unimplemented")
}

// Delete implements entities.AccountStore.
func (m *MemoryAccountsStore) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// GetByID implements entities.AccountStore.
func (m *MemoryAccountsStore) GetByID(ctx context.Context, id string) (entities.Account, error) {
	panic("unimplemented")
}
