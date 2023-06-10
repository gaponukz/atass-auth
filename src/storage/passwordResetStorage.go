package storage

import (
	"auth/src/entities"
	"fmt"
)

type PasswordResetStorage struct {
	storage map[string]entities.GmailWithKeyPair
}

func NewPasswordResetStorage() PasswordResetStorage {
	return PasswordResetStorage{
		storage: make(map[string]entities.GmailWithKeyPair),
	}
}

func (stor PasswordResetStorage) Create(user entities.GmailWithKeyPair) error {
	stor.storage[user.Key] = user
	return nil
}

func (stor PasswordResetStorage) Delete(user entities.GmailWithKeyPair) error {
	delete(stor.storage, user.Key)
	return nil
}

func (stor PasswordResetStorage) GetByUniqueKey(key string) (entities.GmailWithKeyPair, error) {
	user, ok := stor.storage[key]
	if !ok {
		return entities.GmailWithKeyPair{}, fmt.Errorf("user with unique key %s not found", key)
	}
	return user, nil
}
