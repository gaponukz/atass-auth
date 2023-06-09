package storage

import "fmt"

type UserCredentials struct {
	Gmail string `json:"gmail"`
	Key   string `json:"id"`
}

type PasswordResetStorage struct {
	storage map[string]UserCredentials
}

func NewPasswordResetStorage() PasswordResetStorage {
	return PasswordResetStorage{
		storage: make(map[string]UserCredentials),
	}
}

func (stor PasswordResetStorage) Create(user UserCredentials) error {
	stor.storage[user.Key] = user
	return nil
}

func (stor PasswordResetStorage) Delete(user UserCredentials) error {
	delete(stor.storage, user.Key)
	return nil
}

func (stor PasswordResetStorage) GetByUniqueKey(key string) (UserCredentials, error) {
	user, ok := stor.storage[key]
	if !ok {
		return UserCredentials{}, fmt.Errorf("user with unique key %s not found", key)
	}
	return user, nil
}
