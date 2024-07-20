package memory_store

import (
	"github.com/oalexander6/passman/pkg/entities"
)

// MemoryStore is an in-memory store only to be used for testing and mocking.
type MemoryStore struct {
	NotesStore    *MemoryNotesStore
	AccountsStore *MemoryAccountsStore
	Data          *MemoryStoreData
}

type MemoryStoreData struct {
	Notes    []entities.Note
	Accounts []entities.Account
}

func New() *MemoryStore {
	data := &MemoryStoreData{
		Notes:    make([]entities.Note, 0),
		Accounts: make([]entities.Account, 0),
	}

	return &MemoryStore{
		Data:          data,
		NotesStore:    &MemoryNotesStore{Data: data},
		AccountsStore: &MemoryAccountsStore{Data: data},
	}
}
