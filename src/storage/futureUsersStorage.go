package storage

import (
	"auth/src/entities"
	"fmt"
)

type IFutureUserStorage interface {
	Create(entities.FutureUser) error
	Delete(entities.FutureUser) error
	GetByUniqueKey(string) (entities.FutureUser, error)
}

type FutureUserMemoryStorage struct {
	storage map[string]entities.FutureUser
}

func NewFutureUserMemoryStorage() *FutureUserMemoryStorage {
	return &FutureUserMemoryStorage{
		storage: make(map[string]entities.FutureUser),
	}
}

func (stor *FutureUserMemoryStorage) Create(user entities.FutureUser) error {
	stor.storage[user.UniqueKey] = user
	return nil
}

func (stor *FutureUserMemoryStorage) Delete(user entities.FutureUser) error {
	delete(stor.storage, user.UniqueKey)
	return nil
}

func (stor *FutureUserMemoryStorage) GetByUniqueKey(key string) (entities.FutureUser, error) {
	user, ok := stor.storage[key]
	if !ok {
		return entities.FutureUser{}, fmt.Errorf("user with unique key %s not found", key)
	}
	return user, nil
}
