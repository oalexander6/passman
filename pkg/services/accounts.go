package services

import (
	"context"
	"errors"

	"github.com/alexedwards/argon2id"
	"github.com/oalexander6/passman/pkg/entities"
)

func (s *Services) Register(ctx context.Context, accountInput entities.AccountInput) (entities.Account, error) {
	_, err := s.accountsStore.GetByEmail(ctx, accountInput.Email)
	if err == nil {
		return entities.Account{}, entities.ErrAlreadyExists
	}
	if !errors.Is(err, entities.ErrNotFound) {
		return entities.Account{}, err
	}

	hashedPassword, err := argon2id.CreateHash(accountInput.Password, argon2id.DefaultParams)
	if err != nil {
		return entities.Account{}, errors.New("password hash failed")
	}

	accountToStore := entities.AccountInput{
		Email:    accountInput.Email,
		Password: hashedPassword,
	}

	savedAccount, err := s.accountsStore.Create(ctx, accountToStore)
	if err != nil {
		return entities.Account{}, err
	}

	return savedAccount, nil
}

func (s *Services) Login(ctx context.Context, accountLoginInput entities.AccountLoginInput) (entities.Account, error) {
	savedAccount, err := s.accountsStore.GetByEmail(ctx, accountLoginInput.Email)
	if err != nil {
		return entities.Account{}, errors.New("credentials invalid")
	}

	match, err := argon2id.ComparePasswordAndHash(accountLoginInput.Password, savedAccount.Password)
	if err != nil {
		return entities.Account{}, err
	}
	if !match {
		return entities.Account{}, errors.New("credentials invalid")
	}

	return savedAccount, nil
}
